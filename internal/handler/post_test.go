package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"grates/internal/domain"
	"grates/internal/service"
	mock_service "grates/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_createPost(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPost, input domain.Post)

	tests := []struct {
		name                 string
		inputPost            domain.Post
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Valid",
			inputPost: domain.Post{
				Title:   "title",
				Content: "content",
				UsersId: 1,
			},
			inputBody: `{"title":"title","content":"content"}`,
			mockBehavior: func(r *mock_service.MockPost, input domain.Post) {
				r.EXPECT().Create(input).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Invalid input body",
			inputPost:            domain.Post{},
			inputBody:            `{"title":"title"}`,
			mockBehavior:         func(r *mock_service.MockPost, input domain.Post) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name: "Service failure",
			inputPost: domain.Post{
				Title:   "title",
				Content: "content",
				UsersId: 1,
			},
			inputBody: `{"title":"title","content":"content"}`,
			mockBehavior: func(r *mock_service.MockPost, input domain.Post) {
				r.EXPECT().Create(input).Return(0, errors.New("internal creating post error")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal creating post error"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockPost(c)
			test.mockBehavior(repo, test.inputPost)

			services := &service.Service{
				Post: repo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/post", func(c *gin.Context) {
				c.Set(userCtx, domain.User{Id: 1})
			}, handler.createPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/post", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			//assert.Equal(t, test.expectedStatusCode, w.Code)
			//assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
