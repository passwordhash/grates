package handler

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "grates/docs"
	"grates/internal/service"
	"grates/pkg/utils"
	"net/http"
)

type CustomBinding string

const (
	PasswordBinding   CustomBinding = "password"
	LetterWordBinding CustomBinding = "letterword"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes(specialSymbols string) *gin.Engine {
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
		auth.GET("/confirm", h.confirmEmail)
		auth.GET("/check/:email", h.checkEmail)
		auth.POST("/resend/:userId", h.resendEmail)
	}

	api := router.Group("/api", h.userIdentity)
	{
		user := api.Group("/user/:userId")
		{
			user.GET("/posts", h.getUsersPosts)
		}

		profile := api.Group("/profile")
		{
			profile.GET("/:userId", h.getProfile)
			// PROFILE INFO
			profile.PATCH("/", h.updateProfile)
		}

		// FRIENDS
		friends := api.Group("/friends/:userId")
		{
			friends.GET("/", h.friends)
			friends.GET("/requests", h.friendRequests)
			friends.POST("/send-request", h.sendFriendRequest)
			friends.PATCH("/accept", h.acceptFriendRequest)
			friends.PATCH("/unfriend", h.unfriend)
		}

		posts := api.Group("/posts")
		{
			posts.POST("/", h.createPost)
			posts.GET("/:postId", h.getPost)
			// TODO измнить
			posts.GET("/friends/:userId", h.friendsPosts)
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

	if err := h.RegisterNewBindings(specialSymbols); err != nil {
		panic(err)
	}

	return router
}

func (h *Handler) RegisterNewBindings(specialSymbols string) error {
	var isPassword = func(fl validator.FieldLevel) bool {
		if specialSymbols == "" {
			specialSymbols = "!@#$%^&*"
		}
		return utils.IsPassword(fl.Field().String(), specialSymbols)
	}
	var isWord = func(fl validator.FieldLevel) bool {
		return utils.IsOneWordLetter(fl.Field().String())
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation(string(PasswordBinding), isPassword); err != nil {
			return err
		}

		if err := v.RegisterValidation(string(LetterWordBinding), isWord); err != nil {
			return err
		}
	} else {
		return errors.New("binding.Validator.Engine() is not *validator.Validate")
	}

	return nil
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
