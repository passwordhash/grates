package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	if hash == "" {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	err := h.services.Email.ConfirmEmail(hash)
	if errors.Is(err, service.AlreadyConfirmedErr) {
		newResponse(c, http.StatusConflict, "email already confirmed")
		return
	}
	if errors.Is(err, service.HashNotFoundErr) {
		newResponse(c, http.StatusBadRequest, "hash not found")
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal error confirming email")
		return
	}

	logrus.Infof("email confirmed: %s", hash)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

//type resendEmailResponse struct {
//	Hash string `json:"hash"`
//}

// @Summary Resend email
// @Tags auth
// @Description resend confirmation email
// @ID resend-email
// @Accept  json
// @Produce  json
// @Param userId path int true "user's id"
// @Success 200 {object} statusResponse
// @Failure 400,404,500 {object} errorResponse
// @Router /auth/resend/{userId} [post]
func (h *Handler) resendEmail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	err = h.services.Email.ReplaceConfirmationEmail(id)
	if errors.Is(err, service.UserNotFoundError) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, service.AlreadyConfirmedErr) {
		newResponse(c, http.StatusConflict, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal error sending email")
		return
	}

	logrus.Infof("confirmation email was sent to %d", id)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
