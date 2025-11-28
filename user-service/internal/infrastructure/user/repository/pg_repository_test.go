package repository

import (
	"context"
	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-kit/log"
	//"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestUserRepoDBCreate_OK(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	logger := log.NewNopLogger()
	userRepoDB := NewUserRepoDB(logger, sqlxDB)

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
	pendingUser.CreatedAt = time.Now()
	pendingUser.UpdatedAt = pendingUser.CreatedAt
	if err != nil {
		t.Fatal("failed to create user", err)
	}

	columns := []string{"user_id", "first_name", "last_name", "contact", "phone", "email", "password", "date_birthday",
		"gender", "role", "activated", "created_at", "updated_at"}

	rows := sqlmock.NewRows(columns).AddRow(
		pendingUser.UserID.String(),
		pendingUser.FirstName,
		pendingUser.LastName,
		pendingUser.Contact,
		pendingUser.Phone,
		pendingUser.Email,
		pendingUser.Password,
		pendingUser.DateBirthday,
		pendingUser.Gender,
		pendingUser.Role,
		pendingUser.Activated,
		pendingUser.CreatedAt,
		pendingUser.UpdatedAt,
	)

	mock.ExpectQuery(createPendingUserQuery).WithArgs(
		sqlmock.AnyArg(),
		pendingUser.FirstName,
		pendingUser.LastName,
		pendingUser.Contact,
		pendingUser.Phone,
		pendingUser.Email,
		pendingUser.Password,
		pendingUser.DateBirthday,
		pendingUser.Gender,
		pendingUser.Role,
		pendingUser.Activated,
	).WillReturnRows(rows)

	regData := domain.RegisterInput{PersonalDate: &personalData}

	user, err := userRepoDB.Create(context.Background(), &regData)
	require.NoError(t, err)
	require.Equal(t, pendingUser.FirstName, user.FirstName)
	require.Equal(t, pendingUser.LastName, user.LastName)
	require.Equal(t, pendingUser.Contact, user.Contact)
	require.Equal(t, pendingUser.Phone, user.Phone)
	require.Equal(t, pendingUser.Email, user.Email)
	require.Equal(t, pendingUser.Password, user.Password)
	require.Equal(t, pendingUser.DateBirthday, user.DateBirthday)
	require.Equal(t, pendingUser.CreatedAt, user.CreatedAt)
	require.Equal(t, pendingUser.UpdatedAt, user.UpdatedAt)
	require.Equal(t, pendingUser.Role, user.Role)
	require.Equal(t, pendingUser.Activated, user.Activated)
}
