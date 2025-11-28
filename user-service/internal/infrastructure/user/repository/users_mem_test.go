package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-kit/log"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

func TestCreateOK(t *testing.T) {
	logger := log.NewNopLogger()
	userRepoMem := NewUsersDBMem(logger)
	ctx := context.Background()

	personalData := domain.PersonalData{
		FirstName:    "John",
		LastName:     "Doe",
		Contact:      "email",
		Phone:        "+7999999999",
		Email:        "a@b.com",
		Password:     "password",
		DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
	}

	pendingUser, err := domain.NewPendingUser(&personalData)
	if err != nil {
		t.Fatal(err)
	}

	regData := &domain.RegisterInput{&personalData}

	user, err := userRepoMem.Create(ctx, regData)
	if err != nil {
		t.Fatal("create user", err)
	}

	if pendingUser.PersonalData.FirstName != user.PersonalData.FirstName {
		t.Fatalf("FirstName not equal")
	}
	if pendingUser.PersonalData.LastName != user.PersonalData.LastName {
		t.Fatalf("LastName not equal")
	}
	if pendingUser.PersonalData.Contact != user.PersonalData.Contact {
		t.Fatalf("Phone not equal")
	}
	if pendingUser.PersonalData.Phone != user.PersonalData.Phone {
		t.Fatalf("Phone not equal")
	}
	if pendingUser.PersonalData.Email != user.PersonalData.Email {
		t.Fatalf("Email not equal")
	}
	if pendingUser.PersonalData.Password != user.PersonalData.Password {
		t.Fatalf("Password not equal")
	}
	if pendingUser.PersonalData.DateBirthday != user.PersonalData.DateBirthday {
		t.Fatalf("DateBirthday not equal")
	}
	if pendingUser.PersonalData.Gender != user.PersonalData.Gender {
		t.Fatalf("Gender not equal")
	}
	if pendingUser.Activated != false {
		t.Fatalf("Activated should be false")
	}
	if pendingUser.Role != "pending" {
		t.Fatalf("Role not equal")
	}

}

func TestCreateNot_OK(t *testing.T) {
	logger := log.NewNopLogger()
	verifyCodesMemRepo := NewUsersDBMem(logger)
	ctx := context.Background()

	personalData := domain.PersonalData{
		FirstName:    "John",
		LastName:     "Doe",
		Contact:      "email",
		Phone:        "+7999999999",
		Email:        "a@b.com",
		Password:     "password",
		DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
	}

	regData := &domain.RegisterInput{&personalData}

	_, err := verifyCodesMemRepo.Create(ctx, regData)
	if err != nil {
		t.Fatal("create user", err)
	}

	_, err = verifyCodesMemRepo.Create(ctx, regData)
	if err == nil {
		t.Fatal("Must be error")
	}
}

func TestVerifyCodesMem_Concurrency_Insert(t *testing.T) {
	logger := log.NewNopLogger()
	verifyCodesMemRepo := NewUsersDBMem(logger)
	ctx := context.Background()

	goroutines := 10
	perGoroutines := 10

	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := range goroutines {
		id := i
		go func() {
			defer wg.Done()

			for j := range perGoroutines {
				personalData := domain.PersonalData{
					FirstName:    fmt.Sprintf("Joie_%d_%d", id, j),
					LastName:     "Doe",
					Contact:      "email",
					Phone:        "+7999999999",
					Email:        "a@b.com",
					Password:     "password",
					DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
					Gender:       "male",
				}

				regData := &domain.RegisterInput{&personalData}

				_, err := verifyCodesMemRepo.Create(ctx, regData)
				if err != nil {
					t.Fatal("create user", err)
				}

				if _, err := verifyCodesMemRepo.Create(ctx, regData); err != nil {
					t.Error("Error must be nil")
					return
				}
			}
		}()
	}
}
