package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	logrus.Info(header)
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

	user, err := h.services.ParseToken(headerParts[1])

	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, user)
}
