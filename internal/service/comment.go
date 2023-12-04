package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
)

type CommentService struct {
	repo repository.Comment
}

func (c *CommentService) Create(comment domain.CommentCreateInput) (int, error) {
	return c.repo.Create(comment)
}

func (c *CommentService) GetPostComments(postId int) ([]domain.Comment, error) {
	return c.repo.GetPostComments(postId)
}

func (c *CommentService) Delete(userId, commentId int) error {
	return c.repo.Delete(userId, commentId)
}

func (c *CommentService) Update(userId, commentId int, newComment domain.CommentUpdateInput) error {
	return c.repo.Update(userId, commentId, newComment)
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}
