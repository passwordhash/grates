package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
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

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getPost(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPost, postId int, post domain.Post)

	tests := []struct {
		name                 string
		pathParam            string
		postId               int
		post                 domain.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Valid",
			pathParam: "1",
			postId:    1,
			post: domain.Post{
				Id:      1,
				Title:   "title",
				Content: "content",
				UsersId: 1,
				Comments: []domain.Comment{
					{
						Id:      1,
						Content: "content",
					},
				},
				LikesCount:    2,
				CommentsCount: 1,
			},
			mockBehavior: func(r *mock_service.MockPost, postId int, post domain.Post) {
				r.EXPECT().Get(postId).Return(post, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"title":"title","content":"content","users-id":1,"date":"0001-01-01T00:00:00Z","comments":[{"id":1,"content":"content","users-id":0,"posts-id":0,"date":"0001-01-01T00:00:00Z"}],"comments-count":1,"likes-count":2}`,
		},
		{
			name:                 "Invalid path param",
			pathParam:            "abc",
			postId:               0,
			post:                 domain.Post{},
			mockBehavior:         func(r *mock_service.MockPost, postId int, post domain.Post) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid path variable value"}`,
		},
		{
			name:      "Interanl error",
			pathParam: "1",
			postId:    1,
			post:      domain.Post{},
			mockBehavior: func(r *mock_service.MockPost, postId int, post domain.Post) {
				r.EXPECT().Get(postId).Return(domain.Post{}, errors.New("error")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"internal error: error"}`,
		},
		{
			name:      "Post not found",
			postId:    1,
			pathParam: "1",
			post:      domain.Post{},
			mockBehavior: func(r *mock_service.MockPost, postId int, post domain.Post) {
				r.EXPECT().Get(postId).Return(domain.Post{}, service.PostNotFoundErr).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"post not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockPost(c)
			test.mockBehavior(repo, test.postId, test.post)

			services := &service.Service{
				Post: repo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/post/:postId", func(c *gin.Context) {
				c.Set(userCtx, domain.User{Id: 1})
			}, handler.getPost)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/post/"+test.pathParam, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getUsersPosts(t *testing.T) {
	type mockBehavior func(r *mock_service.MockPost, userId int, posts []domain.Post)

	tests := []struct {
		name                 string
		queryParam           string
		posts                []domain.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "Valid",
			queryParam: "userId=1",
			posts: []domain.Post{
				{
					Id:      1,
					Title:   "title",
					Content: "content",
					UsersId: 1,
				},
			},
			mockBehavior: func(r *mock_service.MockPost, userId int, posts []domain.Post) {
				r.EXPECT().GetUsersPosts(userId).Return(posts, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"posts":[{"id":1,"title":"title","content":"content","users-id":1,"date":"0001-01-01T00:00:00Z","comments":null,"likes-count":0}],"count":1}`,
		},
		{
			name: "Valid: 3 posts",
			posts: []domain.Post{
				{
					Id:      1,
					Title:   "title",
					Content: "content",
					UsersId: 1,
				},
				{
					Id:      2,
					Title:   "title",
					Content: "content",
					UsersId: 1,
				},
				{
					Id:      3,
					Title:   "title",
					Content: "content",
					UsersId: 1,
				},
			},
			queryParam: "userId=1",
			mockBehavior: func(r *mock_service.MockPost, userId int, posts []domain.Post) {
				r.EXPECT().GetUsersPosts(userId).Return(posts, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"posts":[{"id":1,"title":"title","content":"content","users-id":1,"date":"0001-01-01T00:00:00Z","comments":null,"likes-count":0},{"id":2,"title":"title","content":"content","users-id":1,"date":"0001-01-01T00:00:00Z","comments":null,"likes-count":0},{"id":3,"title":"title","content":"content","users-id":1,"date":"0001-01-01T00:00:00Z","comments":null,"likes-count":0}],"count":3}`,
		},
		{
			name:       "Valid: post with additionals",
			queryParam: "userId=1",
			posts: []domain.Post{
				{
					Id:      1,
					Title:   "title",
					Content: "content",
					UsersId: 1,
					Comments: []domain.Comment{
						{
							Id:      1,
							Content: "content",
							UsersId: 2,
							PostsId: 1,
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockPost(c)
			test.mockBehavior(repo, 1, test.posts)

			services := &service.Service{
				Post: repo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/posts", func(c *gin.Context) {
				c.Set(userCtx, domain.User{Id: 1})
			}, handler.getUsersPosts)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/posts?"+test.queryParam, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
