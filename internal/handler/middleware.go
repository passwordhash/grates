package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"grates/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		newResponse(c, http.StatusUnauthorized, "auth header is empty")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newResponse(c, http.StatusUnauthorized, "auth header is invalid")
		return
	}

	if headerParts[1] == "" {
		newResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.services.GetUserById(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("error getting user by id: %s", err.Error()))
		return
	}

	// TODO move to another middleware
	// Проверка на подтверждение почты
	if !user.IsConfirmed {
		newResponse(c, http.StatusUnauthorized, "email is not confirmed")
		return
	}

	c.Set(userCtx, user)
}

func (h *Handler) postAffiliation(c *gin.Context) {
	userId := c.MustGet(userCtx).(domain.User).Id

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable data")
		return
	}

	isBelongs, err := h.services.Post.IsPostBelongsToUser(userId, postId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !isBelongs {
		newResponse(c, http.StatusForbidden, fmt.Sprintf("post %s does not belong to the user"))
		return
	}
}
