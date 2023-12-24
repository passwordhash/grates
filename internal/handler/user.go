package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"grates/internal/domain"
	"log"
	"net/http"
	"strconv"
)

type profileResponse struct {
	Profile domain.UserResponse `json:"profile"`
}

// @Summary GetProfile
// @Security ApiKeyAuth
// @Tags profile
// @Description get user profile
// @ID get-profile
// @Accept  json
// @Produce  json
// @Param userId path string true "user id"
// @Success 200 {object} profileResponse
// @Failure 400 {object} statusResponse
// @Router /api/profile/{userId} [get]
func (h *Handler) getProfile(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	profile, err := h.services.GetUserById(userId)
	if err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("can't get user by id: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, profileResponse{
		Profile: domain.UserListToResponse([]domain.User{profile})[0],
	})
}

// @Summary UpdateProfile
// @Security ApiKeyAuth
// @Tags profile
// @Description update user profile
// @ID update-profile
// @Accept  json
// @Produce  json
// @Param input body domain.ProfileUpdateInput true "profile update info"
// @Success 200 {object} statusResponse
// @Failure 400 {object} statusResponse
// @Failure 500 {object} statusResponse
// @Router /api/profile [patch]
func (h *Handler) updateProfile(c *gin.Context) {
	var input domain.ProfileUpdateInput

	userId := c.MustGet("user").(domain.User).Id

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UpdateProfile(userId, input); err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"profile updated"})
}

func (h *Handler) getAllUsers(c *gin.Context) {
	var usersResp []domain.UserResponse

	users, err := h.services.GetAllUsers()
	if err != nil {
		// TEMP
		log.Fatal("get all users error", err.Error())
	}

	for _, user := range users {
		usersResp = append(usersResp, user.ToResponse())
	}

	c.JSON(http.StatusOK, map[string][]domain.UserResponse{
		"users": usersResp,
	})
}
