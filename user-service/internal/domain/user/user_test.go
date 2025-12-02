package user

import (
	"testing"
	"time"
)

func TestNewPendingUser_OK(t *testing.T) {
	pd := PersonalData{
		FirstName:    "John",
		LastName:     "Doe",
		Contact:      "email",
		Phone:        "+7999999999",
		Email:        "a@b.com",
		Password:     "password",
		DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
		Gender:       Male,
	}

	vc, err := NewPendingUser(&pd)
	if err != nil {
		t.Fatal(err)
	}
	if vc.PersonalData.FirstName != pd.FirstName {
		t.Fatalf("FirstName not equal")
	}
	if vc.PersonalData.LastName != pd.LastName {
		t.Fatalf("LastName not equal")
	}
	if vc.PersonalData.Contact != pd.Contact {
		t.Fatalf("Phone not equal")
	}
	if vc.PersonalData.Phone != pd.Phone {
		t.Fatalf("Phone not equal")
	}
	if vc.PersonalData.Email != pd.Email {
		t.Fatalf("Email not equal")
	}
	if vc.PersonalData.Password != pd.Password {
		t.Fatalf("Password not equal")
	}
	if vc.PersonalData.DateBirthday != pd.DateBirthday {
		t.Fatalf("DateBirthday not equal")
	}
	if vc.PersonalData.Gender != pd.Gender {
		t.Fatalf("Gender not equal")
	}
	if vc.Activated != false {
		t.Fatalf("Activated should be false")
	}
	if vc.Role != "pending" {
		t.Fatalf("Role not equal")
	}
}
