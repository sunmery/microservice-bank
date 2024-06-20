// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"context"
)

type Querier interface {
	//CreateAccount
	//
	//  INSERT INTO accounts (owner,
	//                        balance,
	//                        currency)
	//  VALUES ($1,
	//          $2,
	//          $3)
	//  RETURNING id, owner, balance, currency, created_at
	CreateAccount(ctx context.Context, arg *CreateAccountParams) (*Account, error)
	//CreateEntry
	//
	//  INSERT INTO entries(account_id, amount)
	//  VALUES ($1, $2)
	//  RETURNING id, account_id, amount, created_at
	CreateEntry(ctx context.Context, arg *CreateEntryParams) (*Entry, error)
	//CreateTransfer
	//
	//  INSERT INTO transfers(from_account_id, to_account_id, amount)
	//  VALUES ($1, $2, $3)
	//  RETURNING id, from_account_id, to_account_id, amount, created_at
	CreateTransfer(ctx context.Context, arg *CreateTransferParams) (*Transfer, error)
	//DeleteAccount
	//
	//  DELETE
	//  FROM accounts
	//  WHERE id = $1
	//  RETURNING id, owner, balance, currency, created_at
	DeleteAccount(ctx context.Context, id int64) (*Account, error)
	//GetAccount
	//
	//  SELECT id, owner, balance, currency, created_at
	//  FROM accounts
	//  WHERE id = $1
	//  LIMIT 1
	GetAccount(ctx context.Context, id int64) (*Account, error)
	//GetEntry
	//
	//  SELECT id, account_id, amount, created_at
	//  FROM entries
	//  WHERE id = $1
	//  LIMIT 1
	GetEntry(ctx context.Context, id int64) (*Entry, error)
	//GetTransfer
	//
	//  SELECT id, from_account_id, to_account_id, amount, created_at
	//  FROM transfers
	//  WHERE id = $1
	GetTransfer(ctx context.Context, id int64) (*Transfer, error)
	//GetTwoTransfer
	//
	//  SELECT id, from_account_id, to_account_id, amount, created_at
	//  FROM transfers
	//  WHERE from_account_id = $1
	//     OR to_account_id = $2
	//  ORDER BY id
	//  LIMIT $3 OFFSET $4
	GetTwoTransfer(ctx context.Context, arg *GetTwoTransferParams) (*Transfer, error)
	//ListAccounts
	//
	//  SELECT id, owner, balance, currency, created_at
	//  FROM accounts
	//  ORDER BY owner
	//  LIMIT 10
	ListAccounts(ctx context.Context) ([]*Account, error)
	//ListEntries
	//
	//  SELECT id, account_id, amount, created_at
	//  FROM entries
	//  WHERE account_id = $1
	//  ORDER BY id
	//  LIMIT $2 OFFSET $3
	ListEntries(ctx context.Context, arg *ListEntriesParams) (*Entry, error)
	//UpdateAccount
	//
	//  UPDATE accounts
	//  SET owner    = $2,
	//      balance  = $3,
	//      currency = $4
	//  WHERE id = $1
	//  RETURNING id, owner, balance, currency, created_at
	UpdateAccount(ctx context.Context, arg *UpdateAccountParams) (*Account, error)
}

var _ Querier = (*Queries)(nil)
