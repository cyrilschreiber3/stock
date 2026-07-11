-- Advisory lock queries

-- name: CheckAdvisoryExclusiveLock :one
SELECT pg_try_advisory_lock($1::bigint);

-- name: AcquireAdvisoryExclusiveLock :exec
SELECT pg_advisory_lock($1::bigint);

-- name: ReleaseAdvisoryExclusiveLock :one
SELECT pg_advisory_unlock($1::bigint);

-- name: CheckAdvisorySharedLock :one
SELECT pg_try_advisory_lock_shared($1::bigint);

-- name: AcquireAdvisorySharedLock :exec
SELECT pg_advisory_lock_shared($1::bigint);

-- name: ReleaseAdvisorySharedLock :one
SELECT pg_advisory_unlock_shared($1::bigint);

-- Transaction-scoped advisory locks (auto-released at end of transaction, no explicit release needed)

-- name: AcquireAdvisoryXactExclusiveLock :exec
SELECT pg_advisory_xact_lock($1::bigint);

-- name: TryAcquireAdvisoryXactExclusiveLock :one
SELECT pg_try_advisory_xact_lock($1::bigint);

-- name: AcquireAdvisoryXactSharedLock :exec
SELECT pg_advisory_xact_lock_shared($1::bigint);

-- name: TryAcquireAdvisoryXactSharedLock :one
SELECT pg_try_advisory_xact_lock_shared($1::bigint);
