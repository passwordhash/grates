package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"grates/internal/domain"
	"grates/internal/service"
	mock_service "grates/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(r *mock_service.MockUser, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Valid",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().GetUserById(1).Return(domain.User{
					Id:          1,
					IsConfirmed: true,
				}, nil)
				r.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"auth header is empty"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"auth header is invalid"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(r *mock_service.MockUser, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(0, errors.New("invalid token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid token"}`,
		},
		{
			name:        "Get User Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().GetUserById(1).Return(domain.User{}, errors.New("get user error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"error getting user by id: get user error"}`,
		},
		{
			name:        "Email Not Confirmed",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(r *mock_service.MockUser, token string) {
				r.EXPECT().ParseToken(token).Return(1, nil)
				r.EXPECT().GetUserById(1).Return(domain.User{
					Id:          1,
					IsConfirmed: false,
				}, nil)
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"email is not confirmed"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockUser(c)
			test.mockBehavior(repo, test.token)

			services := &service.Service{User: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/identity", handler.userIdentity, func(c *gin.Context) {
				user, ok := c.Get(userCtx)
				if !ok {
					c.String(400, "user not found")
					return
				}
				id := user.(domain.User).Id
				c.String(200, "%d", id)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
