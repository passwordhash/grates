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

func (h *Handler) CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}
