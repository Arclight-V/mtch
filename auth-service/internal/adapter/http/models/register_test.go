package models

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

func TestRegisterRequest_Validator(t *testing.T) {
	tests := []struct {
		name    string
		request *RegisterRequest
		wantErr bool
	}{
		{
			name:    "missing email",
			request: &RegisterRequest{Password: "password", Email: ""},
			wantErr: true,
		},
		{
			name:    "invalid email",
			request: &RegisterRequest{Password: "password", Email: "bad@"},
			wantErr: true,
		},
		{
			name:    "missing password",
			request: &RegisterRequest{Password: "", Email: "example@mail.com"},
			wantErr: true,
		},
		{
			name:    "invalid password",
			request: &RegisterRequest{Password: "123", Email: "example@mail.com"},
			wantErr: true,
		},
		{
			name:    "ok",
			request: &RegisterRequest{Password: "12345678", Email: "example@mail.com"},
			wantErr: false,
		},
	}
	v := validator.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.request)
			if tt.wantErr && err == nil {
				t.Errorf("want error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("want no error, got %v", err)
			}
		})
	}
}
