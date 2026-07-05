package controllers

import (
	"context"

	"github.com/cyrilschreiber3/stock/database"
	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	transactionStateDraft         = "draft"
	transactionStateCompleted     = "completed"
	transactionStatePendingRefund = "pendingRefund"
	transactionTypeBuy            = "buy"
	transactionTypeSell           = "sell"
	transactionTypeAdjustment     = "adjustment"
)

type TransactionController struct {
	db                  *pgxpool.Pool
	queries             *repository.Queries
	inventoryController *InventoryController
}

func NewTransactionController(db *pgxpool.Pool, queries *repository.Queries) *TransactionController {
	return &TransactionController{
		db:                  db,
		queries:             queries,
		inventoryController: NewInventoryController(),
	}
}

func (tc *TransactionController) validateApplyTransition(transactionId uuid.UUID, currentState, targetState string) error {
	if currentState == targetState {
		return &TransactionAlreadyInTargetStateError{
			TransactionID: transactionId.String(),
			TargetState:   currentState,
		}
	}

	if targetState == transactionStateDraft {
		return &TransactionTargetStateNotAllowedError{
			TransactionID: transactionId.String(),
			Method:        "ApplyTransaction",
			CurrentState:  currentState,
			TargetState:   targetState,
			AllowedStates: []string{transactionStateCompleted, transactionStatePendingRefund},
		}
	}

	switch currentState {
	case transactionStateCompleted:
		return &TransactionReadOnlyError{
			TransactionID: transactionId.String(),
			State:         currentState,
		}
	case transactionStatePendingRefund:
		if targetState != transactionStateCompleted {
			return &TransactionReadOnlyError{
				TransactionID: transactionId.String(),
				State:         currentState,
			}
		}
		return nil
	case transactionStateDraft:
		return nil
	default:
		return &TransactionReadOnlyError{
			TransactionID: transactionId.String(),
			State:         currentState,
		}
	}
}

func (tc *TransactionController) ApplyTransaction(ctx context.Context, transactionId uuid.UUID, targetState string) (err error) {
	tx, err := tc.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer database.RollbackTransaction(ctx, tx, &err)

	qtx := tc.queries.WithTx(tx)

	transaction, err := qtx.GetTransactionByID(ctx, transactionId)
	if err != nil {
		return err
	}

	if err := tc.validateApplyTransition(transactionId, transaction.State, targetState); err != nil {
		return err
	}

	if transaction.State == transactionStatePendingRefund && targetState == transactionStateCompleted {
		_, err = qtx.SetTransactionCompleted(ctx, transactionId) // keep same tx
		if err != nil {
			return err
		}
		return tx.Commit(ctx)
	}

	items, err := qtx.GetTransactionItemsByTransactionID(ctx, transactionId)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return &TransactionNoItemsError{
			TransactionID: transactionId.String(),
		}
	}

	switch transaction.TransactionType {
	case transactionTypeBuy:
		for _, item := range items {
			err = tc.inventoryController.IncreaseStockFromBuyItem(ctx, qtx, &item)
			if err != nil {
				return err
			}
		}
	case transactionTypeSell:
		for _, item := range items {
			err = tc.inventoryController.DecreaseStockFromSellItem(ctx, qtx, &item)
			if err != nil {
				return err
			}
		}
	case transactionTypeAdjustment:
		for _, item := range items {
			err = tc.inventoryController.ApplyAdjustmentForItem(ctx, qtx, &item)
			if err != nil {
				return err
			}
		}
	default:
		return &TransactionTypeInvalidError{
			TransactionID:   transactionId.String(),
			TransactionType: transaction.TransactionType,
		}
	}

	switch targetState {
	case transactionStateCompleted:
		_, err = qtx.ApplyTransaction(ctx, transactionId)
		if err != nil {
			return err
		}
	case transactionStatePendingRefund:
		_, err = qtx.ApplyTransactionWithPendingRefund(ctx, transactionId)
		if err != nil {
			return err
		}
	case transactionStateDraft:
		return &TransactionTargetStateNotAllowedError{
			TransactionID: transactionId.String(),
			Method:        "ApplyTransaction",
			CurrentState:  transaction.State,
			TargetState:   targetState,
			AllowedStates: []string{transactionStateCompleted, transactionStatePendingRefund},
		}
	default:
		return &TransactionTargetStateInvalidError{
			TransactionID: transactionId.String(),
			TargetState:   targetState,
		}
	}

	return tx.Commit(ctx)
}
