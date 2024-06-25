package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"user/internal/biz"
	"user/internal/biz/repository"
)

type AccountRepo struct {
	data *Data
	log *log.Helper
}

func NewAccountRepo(
	data *Data,
	logger log.Logger,
	) biz.AccountRepo  {
	return &AccountRepo{
		data: data,
		log: log.NewHelper(logger),
	}
}

func (r *AccountRepo) GetAccount(ctx context.Context, req *biz.GetAccountReq) (*biz.GetAccountRes, error) {
	query := repository.New(r.data.db)
	account, err := query.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &biz.GetAccountRes{
		ID:        account.ID,
		Owner:     account.Owner,
		Balance:   account.Balance,
		Currency:  account.Currency,
	}, nil
}
