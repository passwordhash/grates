package handler

import (
	"github.com/gin-gonic/gin"
	"grates/internal/domain"
	"log"
	"net/http"
)

func (h *Handler) getAllUsers(c *gin.Context) {
	//var users []domain.User

	users, err := h.services.GetAllUsers()
	if err != nil {
		// TEMP
		log.Fatal("get all users error", err.Error())
	}

	c.JSON(http.StatusOK, map[string][]domain.User{
		"users": users,
	})
}
