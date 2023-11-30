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

// @Summary Create
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

	postId, err := h.services.Post.Create(post)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("create post error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"postId": postId,
	})
}

// @Summary Get
// @Security ApiKeyAuth
// @Tags posts
// @Description Get post by id
// @ID get-post
// @Accept json
// @Produce json
// @Param id path int true "post id"
// @Success 200 {object} domain.Post "post info"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId} [get]
func (h *Handler) getPost(c *gin.Context) {
	var post domain.Post

	v := c.Param("postId")

	postId, err := strconv.Atoi(v)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	post, err = h.services.Post.Get(postId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, post)
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
// @Router /api/posts/users/{userId} [patch]
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
		return
	}

	c.JSON(http.StatusOK, usersPostsResponse{
		posts,
		len(posts),
	})
}

// @Summary Update
// @Security ApiKeyAuth
// @Tags posts
// @Description Update post body
// @ID update-post
// @Accept json
// @Produce json
// @Param input body domain.PostUpdateInput true "new post data"
// @Param id path int true "post id"
// @Success 200 {object} statusResponse "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId} [put]
func (h *Handler) updatePost(c *gin.Context) {
	var input domain.PostUpdateInput
	var postId int

	v := c.Param("postId")
	postId, err := strconv.Atoi(v)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable data")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}

	if err := h.services.Post.Update(postId, input); err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("update post error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})

}

// @Sammary Delete
// @Security ApiKeyAuth
// @Tags posts
// @Description Delete post by id
// @ID delete-post
// @Accept json
// @Produce json
// @Param id path int true "post id"
// @Success 200 {string} status "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId} [delete]
func (h *Handler) deletePost(c *gin.Context) {
	v := c.Param("postId")

	id, err := strconv.Atoi(v)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	if err := h.services.Post.Delete(id); err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("delete post error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary CreateComment
// @Security ApiKeyAuth
// @Tags comments
// @Description Create new comment
// @ID create-comment
// @Accept json
// @Produce json
// @Param input body domain.CommentCreateInput true "comment info"
// @Param postId path int true "post id"
// @Success 200 {integer} commentId
// @Failure 400,401,500 {object} errorResponse
// @Router /api/comments/posts/{postId} [post]
func (h *Handler) createComment(c *gin.Context) {
	var commentId int
	var postId int
	var input domain.CommentCreateInput

	// Получение пользователя из контекста
	v, _ := c.Get(userCtx)
	user, ok := v.(domain.User)
	if !ok {
		newResponse(c, http.StatusUnauthorized, "user unauthorized")
		return
	}

	// Получение postId из url
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	// Получение данных комментария из тела запроса
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input")
		return
	}

	input.UserId = user.Id
	input.PostId = postId

	commentId, err = h.services.Comment.Create(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("create comment error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": commentId,
	})
}

type postsCommentsResponse struct {
	Comments []domain.Comment `json:"comments"`
	Count    int              `json:"count"`
}

// @Summary GetPostsComments
// @Security ApiKeyAuth
// @Tags comments
// @Description Get post's comments
// @ID posts-comments
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Success 200 {object} postsCommentsResponse "comments info"
// @Failure 400,500 {object} errorResponse
// @Router /api/comments/posts/{postId} [get]
func (h *Handler) getPostsComments(c *gin.Context) {
	var comments []domain.Comment
	var postId int

	v := c.Param("postId")
	postId, err := strconv.Atoi(v)
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	comments, err = h.services.Comment.GetPostComments(postId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, postsCommentsResponse{
		comments,
		len(comments),
	})
}

// @Summary UpdateComment
// @Security ApiKeyAuth
// @Tags comments
// @Description Update comment body
// @ID update-comment
// @Accept json
// @Produce json
// @Param input body domain.CommentUpdateInput true "new comment data"
// @Param id path int true "comment id"
// TODO
func (h *Handler) updateComment(c *gin.Context) {

}

// @Sammary DeleteComment
// @Security ApiKeyAuth
// @Tags comments
// @Description Delete comment by id
// @ID delete-comment
// @Accept json
// @Produce json
// @Param id path int true "comment id"
// TODO
func (h *Handler) deleteComment(c *gin.Context) {

}
