------------------------------------------------------------------------------------------------------------------------
-- name: CreateState :exec
INSERT INTO state (id, idempotency_key, created_at, updated_at,
                   status, step, type, data, fail_data, meta_data)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetStateByIdempotencyKey :one
SELECT id, idempotency_key, created_at, updated_at,
       status, step, type, data, fail_data, meta_data, error
FROM state
WHERE idempotency_key = $1
LIMIT 1;

-- name: GetStateByID :one
SELECT id, idempotency_key, created_at, updated_at,
       status, step, type, data, fail_data, meta_data, error
FROM state
WHERE id = $1
LIMIT 1;

-- name: UpdateState :one
UPDATE state
SET
    updated_at = $1,
    status = $2,
    step = $3,
    data = $4,
    fail_data = $5,
    meta_data = $6,
    error = $7
WHERE id = $8
RETURNING id;
------------------------------------------------------------------------------------------------------------------------

-- name: SaveStepExecuteInfo :exec
INSERT INTO step_execute_info (
    state_id, start_executed_at, complete_executed_at,
    error, preview_step, next_step
) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetStepExecuteInfos :many
SELECT
    state_id, start_executed_at, complete_executed_at,
    error, preview_step, next_step
FROM step_execute_info
WHERE state_id = $1
ORDER BY start_executed_at;

