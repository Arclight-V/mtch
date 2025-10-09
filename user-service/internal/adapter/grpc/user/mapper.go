package user

// userModeToProto convert models.User to pb.User without passwordHash
//func userModelToProto(user *models.User) *pb.User {
//	out := &pb.User{
//		Uuid:      user.UserID.String(),
//		Email:     user.Email,
//		FirstName: user.FirstName,
//		LastName:  user.LastName,
//		Role:      user.Role,
//		CreatedAt: timestamppb.New(user.CreatedAt),
//		UpdateAt:  timestamppb.New(user.UpdatedAt),
//		Verified:  user.Verified,
//	}
//	if user.Avatar != nil {
//		out.Avatar = *user.Avatar
//	}
//	return out
//}
