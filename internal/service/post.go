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

func (p *PostService) GetWithAdditions(postId int) (domain.Post, error) {
	post, err := p.postRepo.Get(postId)
	if err != nil {
		return post, err
	}

	post.Comments, err = p.commentRepo.GetPostComments(postId)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (p *PostService) GetUsersPosts(userId int) ([]domain.Post, error) {
	posts, err := p.postRepo.GetUsersPosts(userId)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		posts[i].Comments, err = p.commentRepo.GetPostComments(post.Id)
		if err != nil {
			return nil, err
		}
	}
	return posts, nil
}

func (p *PostService) Update(id int, newPost domain.PostUpdateInput) error {
	return p.postRepo.Update(id, newPost)
}

func (p *PostService) Delete(id int) error {
	return p.postRepo.Delete(id)
}
