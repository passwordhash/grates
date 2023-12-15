package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
// @Failure 500 {object} errorResponse
// @Router /auth/confirm/ [get]
func (h *Handler) confirmEmail(c *gin.Context) {
	// TODO: может стоит в запросе передавать еще id пользователя ?
	hash := c.Query("hash")

	err := h.services.Email.ConfirmEmail(hash)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("error confirming email: %s", err.Error()))
		return
	}

	logrus.Infof("email confirmed: %s", hash)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

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

	user, err := h.services.GetUserById(id)
	if err != nil {
		newResponse(c, http.StatusNotFound, fmt.Sprintf("user with id %d not found", id))
		return
	}

	if user.IsConfirmed {
		newResponse(c, http.StatusBadRequest, "email already confirmed")
		return
	}

	err = h.services.Email.ReplaceConfirmationEmail(id, user.Email, user.Name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("error sending email: %s", err.Error()))
		return
	}

	logrus.Infof("confirmation email was sent to %s", user.Email)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
