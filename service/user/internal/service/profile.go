package service

import (
	"context"

	pb "user/api/profile/v1"
)

type ProfileService struct {
	pb.UnimplementedProfileServer
}

func NewProfileService() *ProfileService {
	return &ProfileService{}
}

func (s *ProfileService) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileReply, error) {
	return &pb.CreateProfileReply{}, nil
}
func (s *ProfileService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return &pb.UpdateProfileReply{}, nil
}
func (s *ProfileService) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.DeleteProfileReply, error) {
	return &pb.DeleteProfileReply{}, nil
}
func (s *ProfileService) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileReply, error) {

	return &pb.GetProfileReply{}, nil
}
func (s *ProfileService) ListProfile(ctx context.Context, req *pb.ListProfileRequest) (*pb.ListProfileReply, error) {
	return &pb.ListProfileReply{}, nil
}
