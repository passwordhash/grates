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
// @Security ApiKeyAuth
// @Tags posts
// @Description Create new post
// @ID create-post
// @Accept json
// @Produce json
// @Param input body createPostInput true "post info"
// @Success 200 {integer} postId
// @Failure 400,401,500 {object} errorResponse
// @Router /api/posts [post]
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
// @Security ApiKeyAuth
// @Tags posts
// @Description Get user's posts
// @ID users-posts
// @Accept json
// @Produce json
// @Param userId path int true "user's id"
// @Success 200 {object} usersPostsResponse "post info"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{userId} [get]
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

type updatePostInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// @Summary UpdatePost
// @Param input body updatePostInput true "new post data"
// @Router /api/posts [put]
func (h *Handler) updatePost(c *gin.Context) {
	var input updatePostInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}

}

// @Sammary DeletePost
// @Security ApiKeyAuth
// @Tags posts
// @Description Delete post by id
// @ID delete-post
// @Accept json
// @Produce json
// @Param id path int true "post id"
// @Success 200 {string} status "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{id} [delete]
func (h *Handler) deletePost(c *gin.Context) {
	v := c.Param("id")

	id, err := strconv.Atoi(v)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	if err := h.services.DeletePost(id); err != nil {
		newResponse(c, http.StatusInternalServerError, "post hasn't been deleted")
		return
	}

	c.JSON(http.StatusOK, "ok")
}
