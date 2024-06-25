package service

import (
	"context"
	pb "user/api/account/v1"
	"user/internal/biz"
)

type AccountService struct {
	pb.UnimplementedAccountServiceServer

	pc *biz.AccountUsecase
}

func NewAccountService(pc *biz.AccountUsecase) *AccountService {
	return &AccountService{pc: pc}
}

func (s *AccountService) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReply, error) {
	result, err := s.pc.GetAccount(ctx, &biz.GetAccountReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &pb.GetAccountReply{
		Id: result.ID,
		Name:     result.Owner,
		Owner:    result.Owner,
		Balance:  result.Balance,
		Currency: result.Currency,
	}, nil
}
