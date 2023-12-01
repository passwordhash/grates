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

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refreshTokens)
	}

	api := router.Group("/api", h.userIdentity)
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
		}

		posts := api.Group("/posts")
		{
			posts.POST("/", h.createPost)
			posts.GET("/:postId", h.getPost)
			posts.GET("/user/:userId", h.getUsersPosts)
			posts.PATCH("/:postId", h.updatePost)
			posts.DELETE("/:postId", h.deletePost)

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
