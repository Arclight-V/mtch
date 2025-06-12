package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "proto"
	"user-service/internal/models"
)

func userModelToProto(user *models.User) *pb.User {
	out := &pb.User{
		Uuid:         user.UserID.String(),
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Role:         user.Role,
		PasswordHash: user.PasswordHash,
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdateAt:     timestamppb.New(user.UpdatedAt),
		Verified:     user.Verified,
	}
	if user.Avatar != nil {
		out.Avatar = *user.Avatar
	}
	return out
}
