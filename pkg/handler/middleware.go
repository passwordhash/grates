package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "user"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResp(c, http.StatusUnauthorized, "auth header is empty")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResp(c, http.StatusUnauthorized, "auth header is invalid")
		return
	}

	if headerParts[1] == "" {
		newErrorResp(c, http.StatusUnauthorized, "token is empty")
		return
	}

	user, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResp(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, user)
}
