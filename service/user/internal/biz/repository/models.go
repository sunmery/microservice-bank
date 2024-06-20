// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID        int64              `json:"id"`
	Owner     string             `json:"owner"`
	Balance   int64              `json:"balance"`
	Currency  string             `json:"currency"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type Entry struct {
	ID        int64              `json:"id"`
	AccountID int64              `json:"accountId"`
	Amount    int64              `json:"amount"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

type Transfer struct {
	ID            int64              `json:"id"`
	FromAccountID int64              `json:"fromAccountId"`
	ToAccountID   int64              `json:"toAccountId"`
	Amount        int64              `json:"amount"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt"`
}