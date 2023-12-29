package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"grates/internal/domain"
	"grates/internal/service"
	mock_service "grates/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(userR *mock_service.MockUser, emailR *mock_service.MockEmail, user domain.UserSignUpInput)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.UserSignUpInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid",
			inputBody: `{"email": "test@email.ru", "name":"Test", "password": "qwerty"}`,
			inputUser: domain.UserSignUpInput{
				Name:     "Test",
				Email:    "test@email.ru",
				Password: "qwerty",
			},
			mockBehavior: func(userR *mock_service.MockUser, emailR *mock_service.MockEmail, user domain.UserSignUpInput) {
				userR.EXPECT().CreateUser(user).Return(1417, nil).AnyTimes()
				emailR.EXPECT().ReplaceConfirmationEmail(1417).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1417}`,
		},
		{
			name:                 "Invalid Input",
			inputBody:            `{"email": "testemail.ru", "password": "qwerty}`,
			inputUser:            domain.UserSignUpInput{},
			mockBehavior:         func(userR *mock_service.MockUser, emailR *mock_service.MockEmail, user domain.UserSignUpInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "User With Email Exists",
			inputBody: `{"email": "test@test.ru", "name":"Test", "password": "qwerty"}`,
			inputUser: domain.UserSignUpInput{
				Name:     "Test",
				Email:    "test@test.ru",
				Password: "qwerty",
			},
			mockBehavior: func(userR *mock_service.MockUser, emailR *mock_service.MockEmail, user domain.UserSignUpInput) {
				userR.EXPECT().CreateUser(user).Return(0, service.UserWithEmailExistsError).AnyTimes()
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"message":"user with this email already exists"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_service.NewMockUser(c)
			emailRepo := mock_service.NewMockEmail(c)
			test.mockBehavior(userRepo, emailRepo, test.inputUser)

			services := &service.Service{
				User:  userRepo,
				Email: emailRepo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, input domain.SignInInput)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            domain.SignInInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid",
			inputBody: `{"email": "test@user.ru" , "password": "passwordhash"}`,
			inputUser: domain.SignInInput{
				Email:    "test@user.ru",
				Password: "passwordhash",
			},
			mockBehavior: func(r *mock_service.MockUser, input domain.SignInInput) {
				r.EXPECT().AuthenticateUser(input.Email, input.Password).Return(service.Tokens{
					Access:  "access",
					Refresh: "refresh",
				}, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accessToken":"access","refreshToken":"refresh"}`,
		},
		{
			name:                 "Invalid Input",
			inputBody:            `{"email": "testemail.ru", "pword": "qwerty}`,
			inputUser:            domain.SignInInput{},
			mockBehavior:         func(r *mock_service.MockUser, input domain.SignInInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "User Not Found",
			inputBody: `{"email": "notexists@email.ru",
						 "password": "passwordhash"}`,
			inputUser: domain.SignInInput{
				Email:    "notexists@email.ru",
				Password: "passwordhash",
			},
			mockBehavior: func(r *mock_service.MockUser, input domain.SignInInput) {
				r.EXPECT().AuthenticateUser(input.Email, input.Password).Return(service.Tokens{}, service.UserNotFoundError).AnyTimes()
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth credentials"}`,
		},
		{
			name: "Generate Tokens Error",
			inputBody: `{"email": "test@user.ru",
						 "password": "passwordhash"}`,
			inputUser: domain.SignInInput{
				Email:    "test@user.ru",
				Password: "passwordhash",
			},
			mockBehavior: func(r *mock_service.MockUser, input domain.SignInInput) {
				r.EXPECT().AuthenticateUser(input.Email, input.Password).Return(service.Tokens{}, service.GenerateTokensError).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"error generating tokens"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_service.NewMockUser(c)

			services := &service.Service{
				User: userRepo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			logrus.Info("test.inputBody: ", test.inputBody)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_refresh(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, input string)

	tests := []struct {
		name                 string
		inputBody            string
		inputToken           string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Valid",
			inputBody:  "{\"refreshToken\": \"refreshT\"}",
			inputToken: "refreshT",
			mockBehavior: func(r *mock_service.MockUser, input string) {
				r.EXPECT().RefreshTokens(input).Return(service.Tokens{
					Access:  "accessT",
					Refresh: "refreshT",
				}, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"accessToken":"accessT","refreshToken":"refreshT"}`,
		},
		{
			name:                 "Invalid input body",
			inputBody:            "",
			inputToken:           "",
			mockBehavior:         func(r *mock_service.MockUser, input string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"ivalid input body"}`,
		},
		{
			name:       "Token not found",
			inputBody:  "{\"refreshToken\": \"sometokennotindb\"}",
			inputToken: "sometokennotindb",
			mockBehavior: func(r *mock_service.MockUser, input string) {
				r.EXPECT().RefreshTokens(input).Return(service.Tokens{}, service.RefreshTokenNotFoundError)
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"refresh token not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			userRepo := mock_service.NewMockUser(c)
			test.mockBehavior(userRepo, test.inputToken)

			services := &service.Service{
				User: userRepo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/refresh", handler.refreshTokens)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
