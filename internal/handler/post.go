package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"grates/internal/domain"
	"net/http"
	"strconv"
)

type createPostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// @Summary CreatePost
// @Success 200 {integer} postId
// @Failure 400,401,500 {object} errorResponse
func (h *Handler) createPost(c *gin.Context) {
	var user domain.User
	var post domain.Post

	var input createPostInput

	v, exists := c.Get(userCtx)
	user, ok := v.(domain.User)
	if !ok || !exists {
		newResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	post = domain.Post{
		Title:   input.Title,
		Content: input.Content,
		UsersId: user.Id,
	}

	postId, err := h.services.CreatePost(post)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("create post error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"postId": postId,
	})
}

// @Summary GetPost
func (h *Handler) getPost(c *gin.Context) {

}

type usersPostsResponse struct {
	Posts []domain.Post `json:"posts"`
	Count int           `json:"count"`
}

// @Summary GetUsersPosts
// @Param userId path integer true "user's id"
func (h *Handler) getUsersPosts(c *gin.Context) {
	var posts []domain.Post
	v := c.Param("userId")

	userId, err := strconv.Atoi(v)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	posts, err = h.services.GetUsersPosts(userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, usersPostsResponse{
		posts,
		len(posts),
	})
}
