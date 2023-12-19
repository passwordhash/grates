package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"grates/internal/domain"
	"net/http"
	"strconv"
)

// @Summary SendFriendRequest
// @Tags profile
// @Description send friend request
// @ID send-friend-request
// @Accept  json
// @Produce  json
// @Param toId query string true "user id to send request"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/profile/friend-request/ [post]
func (h *Handler) sendFriendRequest(c *gin.Context) {
	var fromId int
	var toId int

	fromId = c.MustGet(userCtx).(domain.User).Id

	toId, err := strconv.Atoi(c.Query("toId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Friend.SendFriendRequest(fromId, toId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("user %d sent friend request to user %d", fromId, toId)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// @Summary AcceptFriendRequest
// @Tags profile
// @Description accept friend request
// @ID accept-friend-request
// @Accept  json
// @Produce  json
// @Param fromId query string true "user id to accept request"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/profile/accept-request/ [post]
func (h *Handler) acceptFriendRequest(c *gin.Context) {
	var toId int
	var fromId int

	toId = c.MustGet(userCtx).(domain.User).Id

	fromId, err := strconv.Atoi(c.Query("fromId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid query fromId parameter")
		return
	}

	err = h.services.Friend.AcceptFriendRequest(fromId, toId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("user %d accepted friend request from user %d", toId, fromId)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
