package httpadapter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	mem "github.com/ulule/limiter/v3/drivers/store/memory"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	"github.com/Arclight-V/mtch/auth-service/internal/usecase/auth"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/mocks"
)

func TestRegister(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		body string
		stub func(reqUC *mocks.MockRegisterUseCase)
		want want
	}{
		{
			name: "happy-path",
			body: `{
					"first_name":"John", 
					"last_name":"Doe", 
					"contact":"a@b.c", 
					"password":"S3cret42",
					"birth_day":"28", 
					"birth_month":"11", 
					"birth_year":"1992",
					"gender":"male"
					}`,
			stub: func(reqUC *mocks.MockRegisterUseCase) {

				reqUC.
					EXPECT().
					Register(gomock.Any(),
						&auth.RegisterInput{
							FirstName: "John",
							LastName:  "Doe",
							Contact:   "a@b.c",
							Password:  "S3cret42",
							Date:      &auth.Date{BirthDay: 28, BirthMonth: 11, BirthYear: 1992},
						},
					).Return(&auth.RegisterOutput{Email: "a@b.c", Verified: false}, nil).Times(1)
			},
			want: want{http.StatusCreated},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := log.NewNopLogger()
			featuresTest := feature_list.Features{
				feature_list.FeatureKafka:      feature_list.FeatureDisabledByDefault,
				feature_list.VerifyCodeEnabled: feature_list.FeatureEnabledByDefault,
			}
			featureList := feature_list.NewNoopFeatureList(featuresTest)

			regUC := mocks.NewMockRegisterUseCase(ctrl)
			logUC := mocks.NewMockLoginUseCase(ctrl)
			verifyUC := mocks.NewMockVerifyUseCase(ctrl)
			tt.stub(regUC)
			handler := NewHandler(logger, featureList, &Options{}, regUC, logUC, verifyUC)

			req := httptest.NewRequest("POST", apiBase+"auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.router.ServeHTTP(rr, req)

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

			logger := log.NewNopLogger()
			featuresTest := feature_list.Features{
				feature_list.FeatureKafka:      feature_list.FeatureDisabledByDefault,
				feature_list.VerifyCodeEnabled: feature_list.FeatureEnabledByDefault,
			}
			featureList := feature_list.NewNoopFeatureList(featuresTest)

			regUC := mocks.NewMockRegisterUseCase(ctrl)
			logUC := mocks.NewMockLoginUseCase(ctrl)
			verifyUC := mocks.NewMockVerifyUseCase(ctrl)
			handler := NewHandler(logger, featureList, &Options{}, regUC, logUC, verifyUC)

			req := httptest.NewRequest("POST", apiBase+"auth/register", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler.router.ServeHTTP(rr, req)

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

func TestVerifyCode(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name string
		body string
		stub func(reqUC *mocks.MockVerifyUseCase)
		want want
	}{
		{
			name: "happy-code",
			body: `{"code":"123456"}`,
			stub: func(verUC *mocks.MockVerifyUseCase) {

				verUC.
					EXPECT().
					VerifyCode(gomock.Any(), &auth.VerifyInput{Code: "123456"}).
					Return(&auth.VerifyOutput{VerifiedAt: time.Now(), Verified: true}, nil).
					Times(1)
			},
			want: want{http.StatusCreated},
		},
		{
			name: "bad-code",
			body: `{"code":"123456"}`,
			stub: func(verUC *mocks.MockVerifyUseCase) {

				verUC.
					EXPECT().
					VerifyCode(gomock.Any(), &auth.VerifyInput{Code: "123456"}).
					Return(nil, errors.New("any code")).
					Times(1)
			},
			want: want{http.StatusBadRequest},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			logger := log.NewNopLogger()
			featuresTest := feature_list.Features{
				feature_list.FeatureKafka:      feature_list.FeatureDisabledByDefault,
				feature_list.VerifyCodeEnabled: feature_list.FeatureEnabledByDefault,
			}
			featureList := feature_list.NewNoopFeatureList(featuresTest)

			regUC := mocks.NewMockRegisterUseCase(ctrl)
			logUC := mocks.NewMockLoginUseCase(ctrl)
			verifyUC := mocks.NewMockVerifyUseCase(ctrl)
			tt.stub(verifyUC)
			handler := NewHandler(logger, featureList, &Options{}, regUC, logUC, verifyUC)

			req := httptest.NewRequest("POST", apiBase+"auth/verify-code", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.router.ServeHTTP(rr, req)

			require.Equal(t, tt.want.code, rr.Code)

			// TODO: Message?
			//var resp map[string]string
			//_ = json.Unmarshal(rr.Body.Bytes(), &resp)
			//require.Contains(t, resp["message"], tt.want.message)
		})
	}
}
