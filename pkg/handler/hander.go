package handler

import (
	"github.com/gin-gonic/gin"
	"grates/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
	}

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getAllUsers)
		}
	}

	return router
}
