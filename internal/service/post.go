package service

import (
	"fmt"
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

func (s *PostService) Create(post domain.Post) (int, error) {
	return s.postRepo.Create(post)
}

// GetWithAdditions возвращает пост с добавленной информацией.
func (s *PostService) GetWithAdditions(postId int) (domain.Post, error) {
	post, err := s.postRepo.Get(postId)
	if err != nil {
		return post, NotFoundErr{subject: fmt.Sprintf("post with id %d", postId)}
	}

	post.Comments, err = s.commentRepo.GetPostComments(postId)
	if err != nil {
		return post, err
	}

	post.LikesCount, err = s.likeRepo.GetPostLikesCount(postId)
	if err != nil {
		return post, err
	}

	return post, nil
}

// GetUsersPosts возвращает посты пользователя.
func (s *PostService) GetUsersPosts(userId int) ([]domain.Post, error) {

	posts, err := s.postRepo.UsersPosts(userId)
	if err != nil {
		return nil, NotFoundErr{subject: fmt.Sprintf("posts of user with id %d", userId)}
	}

	err = s.fillPostsWithAdditions(&posts)
	if err != nil {
		return nil, fmt.Errorf("fillPostsWithAdditions: %w", err)
	}

	return posts, nil
}

// GetFriendsPosts возвращает посты друзей пользователя.
// Метод сначала получает список id друзей пользователя, затем получает посты по этим id.
// Далее проходится по этим постам и получает комментарии и количество лайков.
func (s *PostService) GetFriendsPosts(userId, limit, offset int) ([]domain.Post, error) {
	var posts []domain.Post

	friendsIds, err := s.frienRepo.FriendUsersIds(userId)
	if err != nil {
		return nil, err
	}

	paramsQuery := fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)

	posts, err = s.postRepo.PostsByUserIds(friendsIds, paramsQuery)
	if err != nil {
		return nil, err
	}

	err = s.fillPostsWithAdditions(&posts)
	if err != nil {
		return nil, fmt.Errorf("fillPostsWithAdditions: %w", err)
	}

	return posts, nil
}

func (s *PostService) Update(id int, newPost domain.PostUpdateInput) error {
	return s.postRepo.Update(id, newPost)
}

func (s *PostService) Delete(id int) error {
	return s.postRepo.Delete(id)
}

// IsPostBelongsToUser проверяет принадлежит ли пользователю пост
func (s *PostService) IsPostBelongsToUser(userId, postId int) (bool, error) {
	post, err := s.postRepo.Get(postId)
	if err != nil {
		return false, err
	}

	return post.UsersId == userId, nil
}

// fillPostsWithAdditions заполняет посты дополнительной информацией.
func (s *PostService) fillPostsWithAdditions(posts *[]domain.Post) error {
	for i, post := range *posts {
		p, err := s.GetWithAdditions(post.Id)
		if err != nil {
			return err
		}

		(*posts)[i] = p
	}

	return nil
}
