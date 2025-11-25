package repository

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

func TestInsertIssueOK(t *testing.T) {
	verifyCodesMemRepo := NewVerifyCodesMem()
	ctx := context.Background()

	vc := domain.VerificationCode{
		UserID:    "uuid",
		Code:      "code",
		Purpose:   domain.EmailVerify,
		ExpiresAt: time.Now().Add(time.Minute * 3),
		Attempts:  0,
	}

	if err := verifyCodesMemRepo.InsertIssue(ctx, &vc); err != nil {
		t.Fatal(err)
	}

}

func TestInsertIssueNot_OK(t *testing.T) {
	verifyCodesMemRepo := NewVerifyCodesMem()
	ctx := context.Background()

	vc := domain.VerificationCode{
		UserID:    "uuid",
		Code:      "code",
		Purpose:   domain.EmailVerify,
		ExpiresAt: time.Now().Add(time.Minute * 3),
		Attempts:  0,
	}

	if err := verifyCodesMemRepo.InsertIssue(ctx, &vc); err != nil {
		t.Fatal("Error must be nil")
	}

	if err := verifyCodesMemRepo.InsertIssue(ctx, &vc); err == nil {
		t.Fatal("Error must not be nil")
	}
}

func TestVerifyCodesMem_Concurrency_Insert(t *testing.T) {
	verifyCodesMemRepo := NewVerifyCodesMem()
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
				vc := domain.VerificationCode{
					UserID:    fmt.Sprintf("uuid_%d_%d", id, j),
					Code:      "code",
					Purpose:   domain.EmailVerify,
					ExpiresAt: time.Now().Add(time.Minute * 3),
					Attempts:  0,
				}
				if err := verifyCodesMemRepo.InsertIssue(ctx, &vc); err != nil {
					t.Error("Error must be nil")
					return
				}
			}
		}()
	}
}
