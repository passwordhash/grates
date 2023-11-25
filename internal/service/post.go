package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(post domain.Post) (int, error) {
	return p.repo.CreatePost(post)
}

func (p *PostService) GetPost(postId int) (domain.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PostService) GetUsersPosts(userId int) ([]domain.Post, error) {
	return p.repo.GetUsersPosts(userId)
}

func (p *PostService) UpdatePost(newPost domain.Post) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostService) DeletePost(id int) error {
	return p.repo.DeletePostById(id)
}
