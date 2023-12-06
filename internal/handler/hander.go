package handler

import (
	"github.com/gin-gonic/gin"
	"grates/internal/service"
	"net/http"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "grates/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/docs/index.html")
	})

	auth := router.Group("/auth", h.CORSMiddleware)
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refreshTokens)
	}

	api := router.Group("/api", h.userIdentity, h.CORSMiddleware)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
		}

		posts := api.Group("/posts")
		{
			posts.POST("/", h.createPost)
			posts.GET("/", h.getUsersPosts)
			posts.GET("/:postId", h.getPost)
			posts.PATCH("/:postId", h.updatePost)
			posts.DELETE("/:postId", h.deletePost)

			posts.POST("/:postId/like", h.likePost)
			posts.DELETE("/:postId/like", h.unlikePost)

			// LIKE: /api/posts/6947/comments
			comments := posts.Group("/:postId/comments")
			{
				comments.POST("/", h.createComment)
				comments.GET("/", h.getPostsComments)
			}

		}

		// QUESTION: ?
		comment := api.Group("/comment")
		{
			comment.PATCH("/:commentId", h.updateComment)
			comment.DELETE("/:commentId", h.deleteComment)
		}
	}

	return router
}
