package database

import (
	"errors"
	"fmt"
)

type AdvisoryLockKey string

const (
	AdvisoryLockKeyMigration AdvisoryLockKey = "migration"
)

var ErrInvalidAdvisoryLockKey = errors.New("invalid advisory lock key")

var advisoryLockKeyRegistry = map[AdvisoryLockKey]int64{
	AdvisoryLockKeyMigration: 1001,
}

func RegisterAdvisoryLockKey(key AdvisoryLockKey, value int64) error {
	if key == "" {
		return fmt.Errorf("%w: empty key name", ErrInvalidAdvisoryLockKey)
	}
	if value == 0 {
		return fmt.Errorf("%w: %s must not be zero", ErrInvalidAdvisoryLockKey, key)
	}

	advisoryLockKeyRegistry[key] = value
	return nil
}

func ResolveAdvisoryLockKey(key AdvisoryLockKey) (int64, error) {
	value, ok := advisoryLockKeyRegistry[key]
	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrInvalidAdvisoryLockKey, key)
	}
	return value, nil
}
