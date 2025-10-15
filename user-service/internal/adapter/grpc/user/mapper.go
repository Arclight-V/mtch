package user

// userModeToProto convert models.User to pb.User without passwordHash
//func userModelToProto(userservice *models.User) *pb.User {
//	out := &pb.User{
//		Uuid:      userservice.UserID.String(),
//		Email:     userservice.Email,
//		FirstName: userservice.FirstName,
//		LastName:  userservice.LastName,
//		Role:      userservice.Role,
//		CreatedAt: timestamppb.New(userservice.CreatedAt),
//		UpdateAt:  timestamppb.New(userservice.UpdatedAt),
//		Verified:  userservice.Verified,
//	}
//	if userservice.Avatar != nil {
//		out.Avatar = *userservice.Avatar
//	}
//	return out
//}
