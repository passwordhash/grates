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
		// Этот запрос должен быть POST, но есть проблема с отправкой формы из письма
		auth.GET("/confirm/", h.confirmEmail)
		auth.POST("/resend/:userId", h.resendEmail)
	}

	api := router.Group("/api", h.userIdentity)
	{
		profile := api.Group("/profile")
		{
			// PROFILE INFO
			profile.PATCH("/", h.updateProfile)
		}

		// FRIENDS
		friends := api.Group("/friends")
		{
			friends.POST("/request", h.sendFriendRequest)
			friends.GET("/:userId", h.friends)
			friends.PATCH("/accept", h.acceptFriendRequest)
			friends.PATCH("/unfriend", h.unfriend)
		}

		posts := api.Group("/posts")
		{
			posts.POST("/", h.createPost)
			posts.GET("/", h.getUsersPosts)
			posts.GET("/friends/:userId", h.friendsPosts)
			posts.GET("/:postId", h.getPost)
			posts.PATCH("/:postId", h.postAffiliation, h.updatePost)
			posts.DELETE("/:postId", h.postAffiliation, h.deletePost)

			posts.POST("/:postId/like", h.likePost)
			posts.DELETE("/:postId/dislike", h.unlikePost)

			// LIKE: /api/posts/6947/comments
			comments := posts.Group("/:postId/comments")
			{
				comments.POST("/", h.createComment)
				comments.GET("/", h.getPostsComments)
			}
		}

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
