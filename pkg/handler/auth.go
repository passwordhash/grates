package handler

import (
	"fmt"
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

	user, err := h.services.GetUserByEmail(input.Email)

	if !user.IsEmtpty() {
		msg := fmt.Sprintf("user with email %s exists", user.Email)
		newErrorResp(c, http.StatusBadRequest, msg)
		return
	}

	id, err := h.services.User.CreateUser(input)
	if err != nil {
		newErrorResp(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
