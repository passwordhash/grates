package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "grates/docs"
	"grates/internal/domain"
	"grates/internal/service"
	"net/http"
	"time"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body SignInInput true "account info"
// @Success 200 {object} idResponse
// @Failure 400,409,500 {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input domain.UserSignUpInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.User.CreateUser(input)
	if errors.Is(err, service.UserWithEmailExistsError) {
		newResponse(c, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		err := h.services.Email.ReplaceConfirmationEmail(id, input.Email, input.Name)
		if err != nil {
			logrus.Errorf("error sending email: %s", err.Error())
			//TODO: подумать над тем, чтобы отправлять письмо повторно
			time.Sleep(5 * time.Second)
			h.services.Email.ReplaceConfirmationEmail(id, input.Email, input.Name)
			return
		}
	}()

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// @Summary SignIn
// @Tags auth
// @Description authenticate account
// @ID login-account
// @Accept       json
// @Produce      json
// @Param input body signInInput true "account credentials"
// @Success      200  {object} signInResponse "tokens"
// @Failure      400,401  {object}  errorResponse
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	var tokens service.Tokens

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	tokens, err := h.services.AuthenticateUser(input.Email, input.Password)
	if errors.Is(err, service.GenerateTokensError) {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if errors.Is(err, service.UserNotFoundError) {
		newResponse(c, http.StatusUnauthorized, fmt.Sprint("invalid auth credentials"))
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	})
}

type refreshInput struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// @Summary RefreshTokens
// @Tags auth
// @Description refresh access and refresh tokens
// @ID refresh-tokens
// @Accept       json
// @Produce      json
// @Param input body refreshInput true "refresh token"
// @Success      200  {object} signInResponse "tokens"
// @Failure      400,401  {object}  errorResponse
// @Router       /auth/refresh [post]
func (h *Handler) refreshTokens(c *gin.Context) {
	var input refreshInput
	if err := c.Bind(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "ivalid input body")
		return
	}

	tokens, err := h.services.RefreshTokens(input.RefreshToken)
	if errors.Is(err, service.GenerateTokensError) {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if errors.Is(err, service.UserNotFoundError) {
		newResponse(c, http.StatusUnauthorized, fmt.Sprintf("refresh token in invalid: %s", err.Error()))
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	})
}
