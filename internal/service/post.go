package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
)

type PostService struct {
	postRepo    repository.Post
	commentRepo repository.Comment
}

func NewPostService(postRepo repository.Post, commentRepo repository.Comment) *PostService {
	return &PostService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (p *PostService) Create(post domain.Post) (int, error) {
	return p.postRepo.Create(post)
}

func (p *PostService) Get(postId int) (domain.Post, error) {
	post, err := p.postRepo.Get(postId)
	if err != nil {
		return domain.Post{}, err
	}

	// TODO: решить как осуществлять проверку ошибки
	comments, err := p.commentRepo.GetPostComments(postId)
	if err != nil {
		return domain.Post{}, err
	}

	post.Comments = comments

	return post, nil
}

func (p *PostService) GetUsersPosts(userId int, commentsLimit int) ([]domain.Post, error) {
	return p.postRepo.GetUsersPosts(userId, commentsLimit)
}

func (p *PostService) Update(id int, newPost domain.PostUpdateInput) error {
	return p.postRepo.Update(id, newPost)
}

func (p *PostService) Delete(id int) error {
	return p.postRepo.Delete(id)
}
