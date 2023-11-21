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

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description authenticate account
// @ID login-account
// @Accept       json
// @Produce      json
// @Param input body signInInput true "account credentials"
// @Success      200  {string} string "token"
// @Failure      400,d401  {object}  errorResponse
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResp(c, http.StatusBadRequest, fmt.Sprintf("bad auth credentials: %s", err.Error()))
		return
	}

	token, err := h.services.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		newErrorResp(c, http.StatusUnauthorized, fmt.Sprintf("invalid auth credentials: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
