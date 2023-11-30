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

func (p *PostService) Create(post domain.Post) (int, error) {
	return p.repo.Create(post)
}

func (p *PostService) Get(postId int) (domain.Post, error) {
	return p.repo.Get(postId)
}

func (p *PostService) GetUsersPosts(userId int) ([]domain.Post, error) {
	return p.repo.GetUsersPosts(userId)
}

func (p *PostService) Update(id int, newPost domain.PostUpdateInput) error {
	return p.repo.Update(id, newPost)
}

func (p *PostService) Delete(id int) error {
	return p.repo.Delete(id)
}
