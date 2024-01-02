package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"grates/internal/domain"
	"grates/internal/service"
	"net/http"
	"strconv"
)

const (
	userIdQuery      = "userId"
	postLimitDefault = 3
)

type createPostInput struct {
	Title   string `json:"title"`
	Content string `json:"content" binding:"required"`
}

// @Summary CreatePost
// @Security ApiKeyAuth
// @Tags posts
// @Description Create new post
// @ID create-post
// @Accept json
// @Produce json
// @Param input body createPostInput true "post info"
// @Success 200 {object} idResponse
// @Failure 400,401,500 {object} errorResponse
// @Router /api/posts [post]
func (h *Handler) createPost(c *gin.Context) {
	var userId int
	var post domain.Post

	var input createPostInput

	userId = c.MustGet(userCtx).(domain.User).Id

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	post = domain.Post{
		Title:   input.Title,
		Content: input.Content,
		UsersId: userId,
	}

	postId, err := h.services.Post.Create(post)
	logrus.Infof("user %d created post %d", userId, postId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "internal creating post error")
		return
	}

	logrus.Infof("user %d created post %d", userId, postId)

	c.JSON(http.StatusOK, idResponse{Id: postId})
}

// @Summary GetPost
// @Security ApiKeyAuth
// @Tags posts
// @Description GetWithAdditions post by id
// @ID get-post
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Success 200 {object} domain.Post "post info"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId} [get]
func (h *Handler) getPost(c *gin.Context) {
	var post domain.Post

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	post, err = h.services.Post.GetWithAdditions(postId)
	if errors.As(err, &service.NotFoundErr{}) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, post)
}

type postsResponse struct {
	Posts []domain.Post `json:"posts"`
	Count int           `json:"count"`
}

// @Summary UsersPosts
// @Security ApiKeyAuth
// @Tags posts
// @Description GetWithAdditions user's posts
// @ID users-posts
// @Accept json
// @Produce json
// @Param userId query int true "user's id"
// @Success 200 {object} postsResponse "post info"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/ [get]
func (h *Handler) getUsersPosts(c *gin.Context) {
	var posts []domain.Post
	var userId int

	userId, err := strconv.Atoi(c.Query(userIdQuery))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid query value of user's id")
		return
	}

	posts, err = h.services.GetUsersPosts(userId)
	if errors.As(err, &service.NotFoundErr{}) {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, postsResponse{
		posts,
		len(posts),
	})
}

// @Summary FriendsPosts
// @Security ApiKeyAuth
// @Tags posts
// @Description GetWithAdditions friends' posts
// @ID friends-posts
// @Accept json
// @Produce json
// @Param userId path int true "user's id"
// @Param limit query int false "limit of posts"
// @Param offset query int false "offset of posts"
// @Success 200 {object} postsResponse "post info"
// @Failure 403,500 {object} errorResponse
// @Router /api/posts/friends/{userId} [get]
func (h *Handler) friendsPosts(c *gin.Context) {
	var userId int
	var limit int
	var offset int

	userId, err := strconv.Atoi(c.Param("userId"))

	if userId != c.MustGet(userCtx).(domain.User).Id {
		newResponse(c, http.StatusForbidden, "you can't get other user's posts")
		return
	}

	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(postLimitDefault)))
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errLimit != nil || errOffset != nil {
		newResponse(c, http.StatusBadRequest, "invalid query value of limit or offset")
		return
	}

	posts, err := h.services.Post.GetFriendsPosts(userId, limit, offset)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, postsResponse{
		posts,
		len(posts),
	})
}

// @Summary UpdatePost
// @Security ApiKeyAuth
// @Tags posts
// @Description Update post body
// @ID update-post
// @Accept json
// @Produce json
// @Param input body domain.PostUpdateInput true "new post data"
// @Param postId path int true "post id"
// @Success 200 {object} statusResponse "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId} [patch]
func (h *Handler) updatePost(c *gin.Context) {
	var input domain.PostUpdateInput
	var postId int

	postId, err := strconv.Atoi(c.Param("postId"))
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

// @Summary DeletePost
// @Security ApiKeyAuth
// @Tags posts
// @Description Delete post by id
// @ID delete-post
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Success 200 {object} statusResponse
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId} [delete]
func (h *Handler) deletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("postId"))
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
// @Success 200 {object} idResponse
// @Failure 400,401,500 {object} errorResponse
// @Router /api/posts/{postId}/comments [post]
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

	c.JSON(http.StatusOK, idResponse{Id: commentId})
}

type postsCommentsResponse struct {
	Comments []domain.Comment `json:"comments"`
	Count    int              `json:"count"`
}

// @Summary GetPostsComments
// @Security ApiKeyAuth
// @Tags comments
// @Description GetWithAdditions post's comments
// @ID posts-comments
// @Accept json
// @Produce json
// @Param postId path int true "post id"
// @Success 200 {object} postsCommentsResponse "comments info"
// @Failure 400,500 {object} errorResponse
// @Router /api/posts/{postId}/comments [get]
func (h *Handler) getPostsComments(c *gin.Context) {
	var comments []domain.Comment
	var postId int

	postId, err := strconv.Atoi(c.Param("postId"))
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
// @Param commentId path int true "comment id"
// @Success 200 {object} statusResponse "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/comment/{commentId} [patch]
func (h *Handler) updateComment(c *gin.Context) {
	var input domain.CommentUpdateInput
	var commentId int

	userId := c.MustGet(userCtx).(domain.User).Id

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable data")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}

	if err := h.services.Comment.Update(userId, commentId, input); err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("update comment error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary DeleteComment
// @Security ApiKeyAuth
// @Tags comments
// @Description Delete comment by id
// @ID delete-comment
// @Accept json
// @Produce json
// @Param commentId path int true "comment id"
// @Success 200 {object} statusResponse "ok"
// @Failure 400,500 {object} errorResponse
// @Router /api/comment/{commentId} [delete]
func (h *Handler) deleteComment(c *gin.Context) {
	var commentId int

	userId := c.MustGet(userCtx).(domain.User).Id

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, "invalid path variable value")
		return
	}

	if err := h.services.Comment.Delete(userId, commentId); err != nil {
		newResponse(c, http.StatusInternalServerError, fmt.Sprintf("delete comment error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func getIntQueryParam(query string) (int, error) {
	if len(query) == 0 {
		return 0, errors.New("empty query")
	}

	return strconv.Atoi(query)
}
