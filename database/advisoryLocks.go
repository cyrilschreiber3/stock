package database

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvisoryLockType string

const (
	AdvisoryLockTypeExclusive AdvisoryLockType = "exclusive"
	AdvisoryLockTypeShared    AdvisoryLockType = "shared"
)

var ErrInvalidAdvisoryLockType = errors.New("invalid advisory lock type")
var ErrNilAdvisoryLockQueries = errors.New("advisory lock queries are nil")

func ParseAdvisoryLockType(value string) (AdvisoryLockType, error) {
	switch AdvisoryLockType(strings.ToLower(strings.TrimSpace(value))) {
	case AdvisoryLockTypeExclusive:
		return AdvisoryLockTypeExclusive, nil
	case AdvisoryLockTypeShared:
		return AdvisoryLockTypeShared, nil
	default:
		return "", fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockType, value)
	}
}

func (lockType AdvisoryLockType) valid() bool {
	switch lockType {
	case AdvisoryLockTypeExclusive, AdvisoryLockTypeShared:
		return true
	default:
		return false
	}
}

type advisoryLockQueries interface {
	CheckAdvisoryExclusiveLock(context.Context, int64) (bool, error)
	AcquireAdvisoryExclusiveLock(context.Context, int64) error
	ReleaseAdvisoryExclusiveLock(context.Context, int64) (bool, error)
	CheckAdvisorySharedLock(context.Context, int64) (bool, error)
	AcquireAdvisorySharedLock(context.Context, int64) error
	ReleaseAdvisorySharedLock(context.Context, int64) (bool, error)
}

func CheckAdvisoryLockInt64(ctx context.Context, queries advisoryLockQueries, key int64, lockType AdvisoryLockType) (bool, error) {
	if !lockType.valid() {
		return false, fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockType, lockType)
	}
	if queries == nil {
		return false, ErrNilAdvisoryLockQueries
	}

	var acquired bool
	var err error
	if lockType == AdvisoryLockTypeShared {
		acquired, err = queries.CheckAdvisorySharedLock(ctx, key)
	} else {
		acquired, err = queries.CheckAdvisoryExclusiveLock(ctx, key)
	}
	if err != nil || !acquired {
		return acquired, err
	}

	_, releaseErr := releaseAdvisoryLockInt64(ctx, queries, key, lockType)
	if releaseErr != nil {
		return false, releaseErr
	}

	return true, nil
}

func acquireAdvisoryLockInt64(ctx context.Context, queries advisoryLockQueries, key int64, lockType AdvisoryLockType) error {
	if !lockType.valid() {
		return fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockType, lockType)
	}
	if queries == nil {
		return ErrNilAdvisoryLockQueries
	}

	if lockType == AdvisoryLockTypeShared {
		return queries.AcquireAdvisorySharedLock(ctx, key)
	}

	return queries.AcquireAdvisoryExclusiveLock(ctx, key)
}

func releaseAdvisoryLockInt64(ctx context.Context, queries advisoryLockQueries, key int64, lockType AdvisoryLockType) (bool, error) {
	if !lockType.valid() {
		return false, fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockType, lockType)
	}
	if queries == nil {
		return false, ErrNilAdvisoryLockQueries
	}

	if lockType == AdvisoryLockTypeShared {
		return queries.ReleaseAdvisorySharedLock(ctx, key)
	}

	return queries.ReleaseAdvisoryExclusiveLock(ctx, key)
}

func CheckAdvisoryLock(ctx context.Context, queries advisoryLockQueries, key AdvisoryLockKey, lockType AdvisoryLockType) (bool, error) {
	resolvedKey, err := ResolveAdvisoryLockKey(key)
	if err != nil {
		return false, err
	}

	return CheckAdvisoryLockInt64(ctx, queries, resolvedKey, lockType)
}

// ---------------------------------------------------------------------------
// Transaction-scoped advisory locks
// ---------------------------------------------------------------------------

// AdvisoryXactLockQueries is satisfied by any tx-bound *repository.Queries.
// These locks are automatically released when the surrounding transaction ends.
type AdvisoryXactLockQueries interface {
	AcquireAdvisoryXactExclusiveLock(context.Context, int64) error
	TryAcquireAdvisoryXactExclusiveLock(context.Context, int64) (bool, error)
	AcquireAdvisoryXactSharedLock(context.Context, int64) error
	TryAcquireAdvisoryXactSharedLock(context.Context, int64) (bool, error)
}

// AcquireAdvisoryXactLock acquires a transaction-scoped advisory lock.
// Must be called with tx-bound queries (qtx := queries.WithTx(tx)).
func AcquireAdvisoryXactLock(ctx context.Context, queries AdvisoryXactLockQueries, key int64, lockType AdvisoryLockType) error {
	if !lockType.valid() {
		return fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockType, lockType)
	}
	if queries == nil {
		return ErrNilAdvisoryLockQueries
	}
	if lockType == AdvisoryLockTypeShared {
		return queries.AcquireAdvisoryXactSharedLock(ctx, key)
	}
	return queries.AcquireAdvisoryXactExclusiveLock(ctx, key)
}

// TryAcquireAdvisoryXactLock tries a transaction-scoped advisory lock without blocking.
// Returns false (no error) if the lock is already held by another session.
func TryAcquireAdvisoryXactLock(ctx context.Context, queries AdvisoryXactLockQueries, key int64, lockType AdvisoryLockType) (bool, error) {
	if !lockType.valid() {
		return false, fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockType, lockType)
	}
	if queries == nil {
		return false, ErrNilAdvisoryLockQueries
	}
	if lockType == AdvisoryLockTypeShared {
		return queries.TryAcquireAdvisoryXactSharedLock(ctx, key)
	}
	return queries.TryAcquireAdvisoryXactExclusiveLock(ctx, key)
}

// AcquireAdvisoryXactLockByKey is the key-name variant of AcquireAdvisoryXactLock.
func AcquireAdvisoryXactLockByKey(ctx context.Context, queries AdvisoryXactLockQueries, key AdvisoryLockKey, lockType AdvisoryLockType) error {
	resolvedKey, err := ResolveAdvisoryLockKey(key)
	if err != nil {
		return err
	}
	return AcquireAdvisoryXactLock(ctx, queries, resolvedKey, lockType)
}

// TryAcquireAdvisoryXactLockByKey is the key-name variant of TryAcquireAdvisoryXactLock.
func TryAcquireAdvisoryXactLockByKey(ctx context.Context, queries AdvisoryXactLockQueries, key AdvisoryLockKey, lockType AdvisoryLockType) (bool, error) {
	resolvedKey, err := ResolveAdvisoryLockKey(key)
	if err != nil {
		return false, err
	}
	return TryAcquireAdvisoryXactLock(ctx, queries, resolvedKey, lockType)
}

// ---------------------------------------------------------------------------
// Session-scoped advisory lock manager
// ---------------------------------------------------------------------------

// SessionLock holds a dedicated pooled connection with an advisory lock on it.
// The lock lives for the lifetime of the connection, not a transaction.
// Always call Release when done or defer it immediately after AcquireSessionLock.
type SessionLock struct {
	conn     *pgxpool.Conn
	queries  *repository.Queries
	key      int64
	lockType AdvisoryLockType
	released bool
}

// AcquireSessionLock pins one connection from pool, acquires a session-scoped
// advisory lock on it, and returns the lock handle.
// The caller must call Release(ctx) when finished.
func AcquireSessionLock(ctx context.Context, pool *pgxpool.Pool, key AdvisoryLockKey, lockType AdvisoryLockType) (*SessionLock, error) {
	resolvedKey, err := ResolveAdvisoryLockKey(key)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("acquiring connection for advisory lock: %w", err)
	}

	q := repository.New(conn)
	if err := acquireAdvisoryLockInt64(ctx, q, resolvedKey, lockType); err != nil {
		conn.Release()
		return nil, fmt.Errorf("acquiring advisory lock: %w", err)
	}

	return &SessionLock{conn: conn, queries: q, key: resolvedKey, lockType: lockType}, nil
}

// Release releases the advisory lock and returns the connection to the pool.
// Safe to call more than once.
func (l *SessionLock) Release(ctx context.Context) error {
	if l.released {
		return nil
	}
	l.released = true
	defer l.conn.Release()

	ok, err := releaseAdvisoryLockInt64(ctx, l.queries, l.key, l.lockType)
	if err != nil {
		return fmt.Errorf("releasing advisory lock: %w", err)
	}
	if !ok {
		return errors.New("advisory lock was not held on this session")
	}
	return nil
}
