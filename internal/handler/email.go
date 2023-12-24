package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"grates/internal/repository"
	"grates/internal/service"
	"net/http"
	"strconv"
)

// @Summary ConfirmEmail
// @Tags auth
// @Description confirm email
// @ID sign-up
// @Accept  json
// @Produce  json
// @Param hash query string true "hash"
// @Success 200 {object} statusResponse
// @Failure 400,409,500 {object} errorResponse
// @Router /auth/confirm/ [get]
func (h *Handler) confirmEmail(c *gin.Context) {
	// TODO: при уже подверждееном email возвращать ошибку
	// TODO: может стоит в запросе передавать еще id пользователя ?
	hash := c.Query("hash")

	err := h.services.Email.ConfirmEmail(hash)
	if errors.Is(err, service.AlreadyConfirmedErr) {
		newResponse(c, http.StatusConflict, fmt.Sprint("email already confirmed"))
		return
	}
	if errors.Is(err, repository.NoChangesErr) {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("invalid hash: %s", hash))
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("error confirming email: %s", err.Error()))
		return
	}

	logrus.Infof("email confirmed: %s", hash)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

type checkEmailResponse struct {
	IsConfirmed bool `json:"is_confirmed" example:"true"`
}

// @Summary CheckEmail
// @Tags auth
// @Description check if user was confirmed by his email
// @ID check-email
// @Accept  json
// @Produce  json
// @Param email path string true "email"
// @Success 200 {object} checkEmailResponse
// @Failure 400 {object} errorResponse
// @Router /auth/check/{email} [get]
func (h *Handler) checkEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := h.services.GetUserByEmail(email)
	if err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("user with email %s not found", email))
		return
	}

	c.JSON(http.StatusOK, checkEmailResponse{IsConfirmed: user.IsConfirmed})
}

type resendEmailResponse struct {
	Hash string `json:"hash"`
}

// @Summary Resend email
// @Tags auth
// @Description resend confirmation email
// @ID resend-email
// @Accept  json
// @Produce  json
// @Param userId path int true "user's id"
// @Success 200 {object} resendEmailResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /auth/resend/{userId} [post]
func (h *Handler) resendEmail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	user, err := h.services.GetUserById(id)
	if err != nil {
		newResponse(c, http.StatusNotFound, fmt.Sprintf("user with id %d not found", id))
		return
	}

	if user.IsConfirmed {
		newResponse(c, http.StatusBadRequest, "email already confirmed")
		return
	}

	hash, err := h.services.Email.ReplaceConfirmationEmail(id, user.Email, user.Name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("error sending email: %s", err.Error()))
		return
	}

	logrus.Infof("confirmation email was sent to %s", user.Email)

	c.JSON(http.StatusOK, resendEmailResponse{Hash: hash})
}
