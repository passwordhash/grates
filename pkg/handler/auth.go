package handler

import (
	"github.com/gin-gonic/gin"
	"grates/internal/entity"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResp(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// TODO: проверка на существование пользователя с такими данными

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
