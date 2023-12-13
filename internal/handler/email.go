package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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

	err = h.services.Email.ReplaceConfirmationEmail(id, user.Email, user.Name)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("error sending email: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
