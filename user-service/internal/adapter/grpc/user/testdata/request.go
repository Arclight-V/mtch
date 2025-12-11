package testdata

import (
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

func NewTestPBRequest() *userservicepb.RegisterRequest {
	return &userservicepb.RegisterRequest{
		PersonalData: &userservicepb.PersonalData{
			FirstName: "John",
			LastName:  "Doe",
			Contact:   "email",
			Phone:     "+7999999999",
			Email:     "a@b.com",
			Password:  "password",
			BirthDate: &userservicepb.Date{BirthDay: 28, BirthMonth: 11, BirthYear: 1992},
			Gender:    "male",
		},
	}
}
