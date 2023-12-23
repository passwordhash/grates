package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"grates/internal/domain"
	"grates/internal/repository"
	"grates/internal/service"
	"net/http"
	"strconv"
)

type friendResponse struct {
	Friends []domain.UserResponse `json:"friends"`
	Count   int                   `json:"count"`
}

// @Summary GetFriends
// @Security ApiKeyAuth
// @Tags friends
// @Description getting user's friends by his id
// @ID get-friends
// @Accept  json
// @Produce  json
// @Param userId path string true "id of user to get friends"
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

type firiendRequestResponse struct {
	Requests []domain.UserResponse `json:"requests"`
	Count    int                   `json:"count"`
}

// @Summary GetFriendRequests
// @Security ApiKeyAuth
// @Tags friends
// @Description getting user's friend requests by his id
// @ID get-friend-requests
// @Accept  json
// @Produce  json
// @Param userId path string true "id of user to get friend requests"
// @Success 200 {object} firiendRequestResponse
// @Failure 400,409,500 {object} errorResponse
// @Router /api/friends/{userId}/requests [get]
func (h *Handler) friendRequests(c *gin.Context) {
	var userId int

	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid user id param")
		return
	}

	if userId != c.MustGet(userCtx).(domain.User).Id {
		newResponse(c, http.StatusForbidden, "you can't get friend requests of another user")
		return
	}

	friends, err := h.services.Friend.FriendRequests(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, firiendRequestResponse{
		Requests: domain.UserListToResponse(friends),
		Count:    len(friends),
	})
}

// @Summary SendFriendRequest
// @Security ApiKeyAuth
// @Tags friends
// @Description send friend request
// @ID send-friend-request
// @Accept  json
// @Produce  json
// @Param userId path string true "user id to send request"
// @Success 200 {object} statusResponse
// @Failure 400,409,500 {object} errorResponse
// @Router /api/friends/{userId}/send-request [post]
func (h *Handler) sendFriendRequest(c *gin.Context) {
	var fromId int
	var toId int

	fromId = c.MustGet(userCtx).(domain.User).Id

	toId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid query toId parameter")
		return
	}

	err = h.services.Friend.SendFriendRequest(fromId, toId)
	if errors.Is(err, service.SelfFriendRequestErr) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if errors.As(err, &repository.CantChangeErr{}) {
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
// @Security ApiKeyAuth
// @Tags friends
// @Description accept friend request
// @ID accept-friend-request
// @Accept  json
// @Produce  json
// @Param userId path string true "user id to accept request"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/friends/{userId}/accept [patch]
func (h *Handler) acceptFriendRequest(c *gin.Context) {
	var fromId int
	var toId int

	toId = c.MustGet(userCtx).(domain.User).Id

	fromId, err := strconv.Atoi(c.Param("userId"))
	logrus.Infof(strconv.Itoa(fromId) + " asdfasdf")
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
// @Security ApiKeyAuth
// @Tags friends
// @Description unfriend user by his id
// @ID unfriend
// @Accept  json
// @Produce  json
// @Param userId path string true "user id to unfriend"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/friends/{userId}/unfriend/ [patch]
func (h *Handler) unfriend(c *gin.Context) {
	var friendId int
	var userId int

	userId = c.MustGet(userCtx).(domain.User).Id

	friendId, err := strconv.Atoi(c.Param("userId"))
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
