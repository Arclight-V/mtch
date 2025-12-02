package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-kit/log"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	user_test_data "github.com/Arclight-V/mtch/user-service/internal/domain/user/testdata"
)

func TestCreateOK(t *testing.T) {
	logger := log.NewNopLogger()
	userRepoMem := NewUsersDBMem(logger)
	ctx := context.Background()

	pendingUser, err := user_test_data.NewTestPendingUser()
	if err != nil {
		t.Fatal(err)
	}

	regData := &domain.RegisterInput{&pendingUser.PersonalData}

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

	pendingUser, err := user_test_data.NewTestPendingUser()
	if err != nil {
		t.Fatal(err)
	}

	regData := &domain.RegisterInput{&pendingUser.PersonalData}

	_, err = verifyCodesMemRepo.Create(ctx, regData)
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

	errCh := make(chan error, goroutines*perGoroutines)
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
					Phone:        fmt.Sprintf("+79999999%d%d", id, j),
					Email:        fmt.Sprintf("a@b%d%d.com", id, j),
					Password:     "password",
					DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
					Gender:       domain.Male,
				}

				regData := &domain.RegisterInput{&personalData}

				if _, errCreate := verifyCodesMemRepo.Create(ctx, regData); errCreate != nil {
					errCh <- errCreate
				}
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for e := range errCh {
		t.Errorf("unexpected error %v", e)
	}

}
