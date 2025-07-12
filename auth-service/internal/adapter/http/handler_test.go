package httpadapter

import (
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	mem "github.com/ulule/limiter/v3/drivers/store/memory"
	goji "goji.io"
	"goji.io/pat"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func routerWithMock(regUC *mocks.MockRegisterUseCase, logUC *mocks.MockLoginUseCase) http.Handler {
	m := goji.NewMux()
	h := NewHandler(regUC, logUC)
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
		stub func(reqUC *mocks.MockRegisterUseCase, logUC *mocks.MockLoginUseCase)
		want want
	}{
		{
			name: "happy-path",
			body: `{"email":"a@b.c", "password":"S3cret42"}`,
			stub: func(reqUC *mocks.MockRegisterUseCase, logUC *mocks.MockLoginUseCase) {
				//passwordHash := "hashed-password"
				//hasher.EXPECT().Hash("S3cret42").Return(passwordHash, nil)
				reqUC.
					EXPECT().
					Register(gomock.Any(), auth.RegisterInput{Email: "a@b.c", Password: "S3cret42"}).
					Return(auth.RegisterOutput{Email: "a@b.c"}, nil).Times(1)
			},
			want: want{http.StatusCreated},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			regUC := mocks.NewMockRegisterUseCase(ctrl)
			logUC := mocks.NewMockLoginUseCase(ctrl)
			tt.stub(regUC, logUC)
			router := routerWithMock(regUC, logUC)

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

			regUC := mocks.NewMockRegisterUseCase(ctrl)
			logUC := mocks.NewMockLoginUseCase(ctrl)
			router := NewRouter(NewHandler(regUC, logUC))

			req := httptest.NewRequest("POST", apiBase+"auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			require.Equal(t, tt.want.code, rr.Code)
		})
	}
}

func TestRateLimiter(t *testing.T) {
	// Limiter
	rate, _ := limiter.NewRateFromFormatted(rateFormatted)
	store := mem.NewStore()
	middleware := mhttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	// fake handler-target
	target := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := middleware.Handler(target)

	for i := int64(1); i <= rate.Limit; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ping", nil)
		req.Header.Set("X-Forwarded-For", "1.2.3.4")

		h.ServeHTTP(rr, req)

		if i <= rate.Limit {
			require.Equal(t, http.StatusOK, rr.Code)
		} else {
			require.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}
}
