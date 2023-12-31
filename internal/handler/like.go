package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"grates/internal/domain"
	"grates/internal/service"
	"net/http"
	"strconv"
)

// @Summary LikePost
// @Security ApiKeyAuth
// @Tags likes
// @Description Like post
// @ID like-post
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Success 200 {object} statusResponse "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId}/like [post]
func (h *Handler) likePost(c *gin.Context) {
	var postId int

	user := c.MustGet(userCtx).(domain.User)

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	err = h.services.Like.LikePost(user.Id, postId)
	if errors.As(err, &service.UserLikePostErr{}) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("like post error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary DislikePost
// @Security ApiKeyAuth
// @Tags likes
// @Description Dislike post
// @ID dislike-post
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Success 200 {object} statusResponse "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId}/dislike [delete]
func (h *Handler) unlikePost(c *gin.Context) {
	var postId int

	user := c.MustGet(userCtx).(domain.User)

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	err = h.services.Like.UnlikePost(user.Id, postId)
	if errors.As(err, &service.UserLikePostErr{}) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("unlike post error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
