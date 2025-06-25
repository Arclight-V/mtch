package httpadapter

import (
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
	"net/http/httptest"
	pb "proto"
	"strings"
	"testing"
)

func routerWithMock(uc *mocks.MockUserRepo, ts *mocks.MockTokenSigner) http.Handler {
	m := goji.NewMux()
	i := auth.Interactor{UserRepo: uc, TokenSigner: ts}
	h := NewHandler(&i)
	m.HandleFunc(pat.Post(apiBase+"auth/register"), h.Register)
	return m
}

func TestRegister(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		body string
		stub func(*mocks.MockUserRepo, *mocks.MockTokenSigner)
		want want
	}{
		{
			name: "happy-path",
			body: `{"email":"a@b.c", "password":"S3cret42"}`,
			stub: func(uc *mocks.MockUserRepo, ts *mocks.MockTokenSigner) {
				uc.
					EXPECT().
					Register(gomock.Any(), &pb.RegisterRequest{Email: "a@b.c", Password: "S3cret42"}).
					Return(&pb.RegisterResponse{User: &pb.User{Email: "a@b.c", PasswordHash: "S3cret42"}}, nil).
					Times(1)
			},
			want: want{http.StatusCreated},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockUserRepo(ctrl)
			ts := mocks.NewMockTokenSigner(ctrl)
			tt.stub(uc, ts)
			router := routerWithMock(uc, ts)

			req := httptest.NewRequest("POST", apiBase+"auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, tt.want.code, rr.Code)

			// TODO: Message?
			//var resp map[string]string
			//_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			//require.Contains(t, resp["message"], tt.want.message)
		})
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	type want struct {
		code int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "missing email",
			body: `{"password":""}`,
			want: want{http.StatusBadRequest},
		},

		{
			name: "invalid email",
			body: `{"password":"bad@"}`,
			want: want{http.StatusBadRequest},
		},

		{
			name: "missing password",
			body: `{"":"example@mail.com"}`,
			want: want{http.StatusBadRequest},
		},

		{
			name: "invalid password",
			body: `{"123":"example@mail.com"}`,
			want: want{http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uc := mocks.NewMockUserRepo(ctrl)
			ts := mocks.NewMockTokenSigner(ctrl)
			router := NewRouter(NewHandler(&auth.Interactor{UserRepo: uc, TokenSigner: ts}))

			req := httptest.NewRequest("POST", apiBase+"auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			require.Equal(t, tt.want.code, rr.Code)
		})
	}
}
