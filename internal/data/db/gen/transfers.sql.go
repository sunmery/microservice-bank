// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transfers.sql

package gen

import (
	"context"
)

const CreateTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers(from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING id, from_account_id, to_account_id, amount, created_at
`

type CreateTransferParams struct {
	FromAccountID int64 `json:"fromAccountId"`
	ToAccountID   int64 `json:"toAccountId"`
	Amount        int64 `json:"amount"`
}

// CreateTransfer
//
//	INSERT INTO transfers(from_account_id, to_account_id, amount)
//	VALUES ($1, $2, $3)
//	RETURNING id, from_account_id, to_account_id, amount, created_at
func (q *Queries) CreateTransfer(ctx context.Context, arg *CreateTransferParams) (*Transfer, error) {
	row := q.db.QueryRow(ctx, CreateTransfer, arg.FromAccountID, arg.ToAccountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return &i, err
}

const GetTransfer = `-- name: GetTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
WHERE id = $1
`

// GetTransfer
//
//	SELECT id, from_account_id, to_account_id, amount, created_at
//	FROM transfers
//	WHERE id = $1
func (q *Queries) GetTransfer(ctx context.Context, id int64) (*Transfer, error) {
	row := q.db.QueryRow(ctx, GetTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return &i, err
}

const GetTwoTransfer = `-- name: GetTwoTransfer :one
SELECT id, from_account_id, to_account_id, amount, created_at
FROM transfers
WHERE from_account_id = $1
   OR to_account_id = $2
ORDER BY id
LIMIT $3 OFFSET $4
`

type GetTwoTransferParams struct {
	FromAccountID int64 `json:"fromAccountId"`
	ToAccountID   int64 `json:"toAccountId"`
	Limit         int64 `json:"limit"`
	Offset        int64 `json:"offset"`
}

// GetTwoTransfer
//
//	SELECT id, from_account_id, to_account_id, amount, created_at
//	FROM transfers
//	WHERE from_account_id = $1
//	   OR to_account_id = $2
//	ORDER BY id
//	LIMIT $3 OFFSET $4
func (q *Queries) GetTwoTransfer(ctx context.Context, arg *GetTwoTransferParams) (*Transfer, error) {
	row := q.db.QueryRow(ctx, GetTwoTransfer,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return &i, err
}
