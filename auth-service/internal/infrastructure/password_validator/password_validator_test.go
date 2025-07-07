package password_validator

import "testing"

func Test_InvalidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "21 characters",
			password: "abc12345asb44444444444",
			wantErr:  true,
		},
		{
			name:     "unvalidated characters",
			password: "невалидныйпароль",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		passwordValidator := NewUserPasswordValidator()
		t.Run(tt.name, func(t *testing.T) {
			err := passwordValidator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_ValidPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "abc12345Dasb!44",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		passwordValidator := NewUserPasswordValidator()
		t.Run(tt.name, func(t *testing.T) {
			err := passwordValidator.Validate(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
