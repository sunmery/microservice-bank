-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
    FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY owner
LIMIT 10;

-- name: UpdateAccount :one
UPDATE accounts
SET owner    = $2,
    balance  = $3,
    currency = $4
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :one
DELETE
FROM accounts
WHERE id = $1
RETURNING *;

-- name: CreateAccount :one
INSERT INTO accounts (owner,
                      balance,
                      currency)
VALUES ($1,
        $2,
        $3)
RETURNING *;
