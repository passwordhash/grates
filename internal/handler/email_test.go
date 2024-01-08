package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"grates/internal/service"
	mock_service "grates/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_confirm(t *testing.T) {
	type mockBehavior func(r *mock_service.MockEmail, input string)

	tests := []struct {
		name                 string
		queryHash            string
		inputHash            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid",
			queryHash: "hash=somehashforconfirm",
			inputHash: "somehashforconfirm",
			mockBehavior: func(r *mock_service.MockEmail, input string) {
				r.EXPECT().ConfirmEmail("somehashforconfirm").Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:                 "Invalid input",
			queryHash:            "",
			inputHash:            "",
			mockBehavior:         func(r *mock_service.MockEmail, input string) {},
			expectedResponseBody: `{"message":"invalid input body"}`,
			expectedStatusCode:   400,
		},
		{
			name:      "Already confirmed",
			queryHash: "hash=somehash",
			inputHash: "somehash",
			mockBehavior: func(r *mock_service.MockEmail, input string) {
				r.EXPECT().ConfirmEmail("somehash").Return(service.AlreadyConfirmedErr)
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"email already confirmed"}`,
		},
		{
			name:      "Hash not found",
			queryHash: "hash=somehash",
			inputHash: "somehash",
			mockBehavior: func(r *mock_service.MockEmail, input string) {
				r.EXPECT().ConfirmEmail("somehash").Return(service.HashNotFoundErr)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"hash not found"}`,
		},
		{
			name:      "Internal error",
			queryHash: "hash=somehash",
			inputHash: "somehash",
			mockBehavior: func(r *mock_service.MockEmail, input string) {
				r.EXPECT().ConfirmEmail("somehash").Return(errors.New("some internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal error confirming email"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			emailRepo := mock_service.NewMockEmail(c)
			test.mockBehavior(emailRepo, test.inputHash)

			services := &service.Service{
				Email: emailRepo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/confirm", handler.confirmEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/confirm?%s", test.queryHash), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_resend(t *testing.T) {
	type mockBehavior func(r *mock_service.MockEmail, input int)

	tests := []struct {
		name                 string
		queryPath            string
		inputId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid",
			queryPath: "1",
			inputId:   1,
			mockBehavior: func(r *mock_service.MockEmail, input int) {
				r.EXPECT().ReplaceConfirmationEmail(input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"ok"}`,
		},
		{
			name:      "Invalid path var",
			queryPath: "asdg",
			inputId:   1,
			mockBehavior: func(r *mock_service.MockEmail, input int) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid path variable value"}`,
		},
		{
			name:      "Bad user id",
			queryPath: "1",
			inputId:   1,
			mockBehavior: func(r *mock_service.MockEmail, input int) {
				r.EXPECT().ReplaceConfirmationEmail(input).Return(service.UserNotFoundError)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"user not found"}`,
		},
		{
			name:      "Email already confirmed",
			queryPath: "1",
			inputId:   1,
			mockBehavior: func(r *mock_service.MockEmail, input int) {
				r.EXPECT().ReplaceConfirmationEmail(input).Return(service.AlreadyConfirmedErr)
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"email already confirmed"}`,
		},
		{
			name:      "Interanl error",
			queryPath: "1",
			inputId:   1,
			mockBehavior: func(r *mock_service.MockEmail, input int) {
				r.EXPECT().ReplaceConfirmationEmail(input).Return(errors.New("internal error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal error sending email"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			emailRepo := mock_service.NewMockEmail(c)
			test.mockBehavior(emailRepo, test.inputId)

			services := &service.Service{
				Email: emailRepo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/resend/:userId", handler.resendEmail)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", fmt.Sprintf("/resend/%s", test.queryPath), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
