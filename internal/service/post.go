package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
)

type PostService struct {
	postRepo    repository.Post
	commentRepo repository.Comment
	likeRepo    repository.Like
	frienRepo   repository.Friend
}

func NewPostService(postRepo repository.Post, commentRepo repository.Comment, like repository.Like, friend repository.Friend) *PostService {
	return &PostService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		likeRepo:    like,
		frienRepo:   friend,
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

	post.LikesCount, err = p.likeRepo.GetPostLikesCount(postId)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (p *PostService) GetUsersPosts(userId int) ([]domain.Post, error) {
	posts, err := p.postRepo.UsersPosts(userId)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		posts[i].Comments, err = p.commentRepo.GetPostComments(post.Id)
		if err != nil {
			return nil, err
		}

		posts[i].LikesCount, err = p.likeRepo.GetPostLikesCount(post.Id)
		if err != nil {
			return nil, err
		}
	}

	return posts, nil
}

// GetFriendsPosts возвращает посты друзей пользователя.
// Метод сначала получает список id друзей пользователя, затем получает посты по этим id.
// Далее проходится по этим постам и получает комментарии и количество лайков.
func (p *PostService) GetFriendsPosts(userId int) ([]domain.Post, error) {
	var posts []domain.Post

	friendsIds, err := p.frienRepo.FriendUsersIds(userId)
	if err != nil {
		return nil, err
	}

	posts, err = p.postRepo.PostsByUserIds(friendsIds)
	if err != nil {
		return nil, err
	}

	for i, post := range posts {
		posts[i].Comments, err = p.commentRepo.GetPostComments(post.Id)
		if err != nil {
			return nil, err
		}

		posts[i].LikesCount, err = p.likeRepo.GetPostLikesCount(post.Id)
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

// IsPostBelongsToUser проверяет принадлежит ли пользователю пост
func (p *PostService) IsPostBelongsToUser(userId, postId int) (bool, error) {
	post, err := p.postRepo.Get(postId)
	if err != nil {
		return false, err
	}

	return post.UsersId == userId, nil
}
