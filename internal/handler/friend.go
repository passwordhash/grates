package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"grates/internal/domain"
	"grates/internal/service"
	"net/http"
	"strconv"
)

type friendResponse struct {
	Friends []domain.UserResponse `json:"friends"`
	Count   int                   `json:"count"`
}

// @Summary GetFriends
// @Tags profile
// @Description get friends
// @ID get-friends
// @Accept  json
// @Produce  json
// @Param userId path string true "user id"
// @Success 200 {object} friendResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/friends/{userId} [get]
func (h *Handler) friends(c *gin.Context) {
	var userId int

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}

	friends, err := h.services.Friend.GetFriends(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, friendResponse{
		Friends: domain.UserListToResponse(friends),
		Count:   len(friends),
	})
}

// @Summary SendFriendRequest
// @Tags profile
// @Description send friend request
// @ID send-friend-request
// @Accept  json
// @Produce  json
// @Param toId query string true "user id to send request"
// @Success 200 {object} statusResponse
// @Failure 400,409,500 {object} errorResponse
// @Router /api/friends/request [post]
func (h *Handler) sendFriendRequest(c *gin.Context) {
	var fromId int
	var toId int

	fromId = c.MustGet(userCtx).(domain.User).Id

	toId, err := strconv.Atoi(c.Query("toId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid query toId parameter")
		return
	}

	err = h.services.Friend.SendFriendRequest(fromId, toId)
	if errors.Is(err, service.SelfFriendRequestErr) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, service.AleadySendErr) {
		newResponse(c, http.StatusConflict, err.Error())
		return
	}
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
// @Router /api/friends/accept [patch]
func (h *Handler) acceptFriendRequest(c *gin.Context) {
	var fromId int
	var toId int

	toId = c.MustGet(userCtx).(domain.User).Id

	fromId, err := strconv.Atoi(c.Query("fromId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid query fromId parameter")
		return
	}

	err = h.services.Friend.AcceptFriendRequest(fromId, toId)
	if err != nil && !errors.As(err, &service.InternalErr{}) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("user %d accepted friend request from user %d", toId, fromId)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

// @Summary Unfriend
// @Tags profile
// @Description unfriend
// @ID unfriend
// @Accept  json
// @Produce  json
// @Param friendId query string true "user id to unfriend"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/friends/unfriend [patch]
func (h *Handler) unfriend(c *gin.Context) {
	var friendId int
	var userId int

	userId = c.MustGet(userCtx).(domain.User).Id

	friendId, err := strconv.Atoi(c.Query("friendId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid query friendId parameter")
		return
	}

	err = h.services.Friend.Unfriend(userId, friendId)
	if errors.Is(err, service.SelfFriendRequestErr) || errors.As(err, &service.NotFoundErr{}) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("user %d unfriended user %d", userId, friendId)

	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
