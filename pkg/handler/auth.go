package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "grates/docs"
	"grates/internal/entity"
	"net/http"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body entity.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Router /auth/sign-up [post]
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
