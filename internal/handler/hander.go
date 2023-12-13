package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "grates/docs"
	"grates/internal/service"
	"net/http"
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

	router.Use(cors.Default())
	router.Use(corsMiddleware())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refreshTokens)
		//auth.POST("/confirm", h.confirmEmail)
		auth.POST("/resend/:userId", h.resendEmail)
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

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
