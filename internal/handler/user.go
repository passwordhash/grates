package handler

import (
	"github.com/gin-gonic/gin"
	"grates/internal/domain"
	"log"
	"net/http"
)

// @Summary UpdateProfile
// @Security ApiKeyAuth
// @Tags user
// @Description update user profile
// @ID update-profile
// @Accept  json
// @Produce  json
// @Param input body domain.ProfileUpdateInput true "profile update info"
// @Success 200 {object} statusResponse
// @Failure 400 {object} statusResponse
// @Failure 500 {object} statusResponse
// @Router /user/profile [patch]
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
