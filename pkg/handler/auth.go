package handler

import (
	"github.com/gin-gonic/gin"
	"grates/internal/entity"
	"log"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.User

	if err := c.BindJSON(&input); err != nil {
		// TODO: custom response
		log.Fatal("Json binidng error", err.Error())
		return
	}

	// TODO: проверка на существование пользователя с такими данными

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		// TODO: custom error
		log.Fatal(err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
