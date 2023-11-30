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

func (c *CommentService) Delete(id int) error {
	//TODO implement me
	panic("implement me")
}

func (c *CommentService) Update(id int, newComment domain.CommentCreateInput) error {
	//TODO implement me
	panic("implement me")
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}
