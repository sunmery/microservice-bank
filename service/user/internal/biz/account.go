package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetAccountReq struct {
	Id    int64
}

// type GetAccountRes repository.Account
type GetAccountRes struct {
	ID        int64              `json:"id"`
	Owner     string             `json:"owner"`
	Balance   int64              `json:"balance"`
	Currency  string             `json:"currency"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}
// type CreateAccountReq repository.CreateAccountParams

type AccountRepo interface {
	// CreateAccount
	//
	//  INSERT INTO accounts (owner,
	//                        balance,
	//                        currency)
	//  VALUES ($1,
	//          $2,
	//          $3)
	//  RETURNING id, owner, balance, currency, created_at
	// CreateAccount(ctx context.Context, arg *CreateAccountReq) (*GetAccountRes, error)
	// DeleteAccount
	//
	//  DELETE
	//  FROM accounts
	//  WHERE id = $1
	//  RETURNING id, owner, balance, currency, created_at
	// DeleteAccount(ctx context.Context, id int64) (*GetAccountRes, error)
	// GetAccount
	//
	//  SELECT id, owner, balance, currency, created_at
	//  FROM accounts
	//  WHERE id = $1
	//  LIMIT 1
	GetAccount(ctx context.Context, req *GetAccountReq) (*GetAccountRes, error)
	// GetAccountForUpdate
	//
	//  SELECT id, owner, balance, currency, created_at
	//  FROM accounts
	//  WHERE id = $1
	//      FOR NO KEY UPDATE
	// GetAccountForUpdate(ctx context.Context, id int64) (*GetAccountRes, error)
}

type AccountUsecase struct {
	repo AccountRepo
	log *log.Helper
}

func NewAccountUsecase(repo AccountRepo, logger log.Logger) *AccountUsecase  {
	return &AccountUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (ac *AccountUsecase) GetAccount(ctx context.Context, req *GetAccountReq) (*GetAccountRes, error) {
	ac.log.WithContext(ctx).Infof("GetAccount: %v", req)
	return ac.repo.GetAccount(ctx, req)
}
