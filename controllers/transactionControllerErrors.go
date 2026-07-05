package controllers

import (
	"fmt"
	"net/http"
)

type TransactionNotFoundError struct {
	TransactionID string
}

func (e *TransactionNotFoundError) Error() string {
	return "transaction not found: " + e.TransactionID
}

func (e *TransactionNotFoundError) HumanReadableError() string {
	return "Transaction not found"
}

func (e *TransactionNotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

type TransactionReadOnlyError struct {
	TransactionID string
	State         string
}

func (e *TransactionReadOnlyError) Error() string {
	return fmt.Sprintf("transaction %s in state %s is not writable", e.TransactionID, e.State)
}

func (e *TransactionReadOnlyError) HumanReadableError() string {
	return fmt.Sprintf("Transaction cannot be modified in state %s", e.State)
}

func (e *TransactionReadOnlyError) HTTPStatusCode() int {
	return http.StatusConflict
}

type TransactionTypeInvalidError struct {
	TransactionID   string
	TransactionType string
}

func (e *TransactionTypeInvalidError) Error() string {
	return fmt.Sprintf("transaction %s has invalid type %s", e.TransactionID, e.TransactionType)
}

func (e *TransactionTypeInvalidError) HumanReadableError() string {
	return fmt.Sprintf("Invalid transaction type: %s", e.TransactionType)
}

func (e *TransactionTypeInvalidError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type TransactionNoItemsError struct {
	TransactionID string
}

func (e *TransactionNoItemsError) Error() string {
	return fmt.Sprintf("transaction %s has no items", e.TransactionID)
}

func (e *TransactionNoItemsError) HumanReadableError() string {
	return "Transaction must contain at least one item"
}

func (e *TransactionNoItemsError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type TransactionStateInvalidError struct {
	TransactionID string
	State         string
}

func (e *TransactionStateInvalidError) Error() string {
	return fmt.Sprintf("transaction %s has invalid state %s", e.TransactionID, e.State)
}

func (e *TransactionStateInvalidError) HumanReadableError() string {
	return fmt.Sprintf("Invalid transaction state: %s", e.State)
}

func (e *TransactionStateInvalidError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type TransactionStateNotAllowedError struct {
	TransactionID string
	Method        string
	CurrentState  string
	AllowedStates []string
}

func (e *TransactionStateNotAllowedError) Error() string {
	return fmt.Sprintf("method %s cannot run for transaction %s in state %s. Allowed states: %v", e.Method, e.TransactionID, e.CurrentState, e.AllowedStates)
}

func (e *TransactionStateNotAllowedError) HumanReadableError() string {
	return fmt.Sprintf("Operation not allowed while transaction is in state %s", e.CurrentState)
}

func (e *TransactionStateNotAllowedError) HTTPStatusCode() int {
	return http.StatusConflict
}

type TransactionTargetStateInvalidError struct {
	TransactionID string
	TargetState   string
}

func (e *TransactionTargetStateInvalidError) Error() string {
	return fmt.Sprintf("transaction %s cannot transition to invalid target state %s", e.TransactionID, e.TargetState)
}

func (e *TransactionTargetStateInvalidError) HumanReadableError() string {
	return fmt.Sprintf("Invalid target state: %s", e.TargetState)
}

func (e *TransactionTargetStateInvalidError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type TransactionTargetStateNotAllowedError struct {
	TransactionID string
	Method        string
	CurrentState  string
	TargetState   string
	AllowedStates []string
}

func (e *TransactionTargetStateNotAllowedError) Error() string {
	return fmt.Sprintf("method %s cannot transition transaction %s from state %s to target state %s. Allowed target states: %v", e.Method, e.TransactionID, e.CurrentState, e.TargetState, e.AllowedStates)
}

func (e *TransactionTargetStateNotAllowedError) HumanReadableError() string {
	return fmt.Sprintf("Transition from %s to %s is not allowed", e.CurrentState, e.TargetState)
}

func (e *TransactionTargetStateNotAllowedError) HTTPStatusCode() int {
	return http.StatusConflict
}

type TransactionAlreadyInTargetStateError struct {
	TransactionID string
	TargetState   string
}

func (e *TransactionAlreadyInTargetStateError) Error() string {
	return fmt.Sprintf("transaction %s is already in target state %s", e.TransactionID, e.TargetState)
}

func (e *TransactionAlreadyInTargetStateError) HumanReadableError() string {
	return fmt.Sprintf("Transaction is already in state %s", e.TargetState)
}

func (e *TransactionAlreadyInTargetStateError) HTTPStatusCode() int {
	return http.StatusConflict
}
