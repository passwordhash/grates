package service

import (
	"fmt"
	"grates/internal/repository"
)

type UserLikePostErr struct {
	msg string
}

func (l UserLikePostErr) Error() string {
	return l.msg
}

type LikeService struct {
	likeRepo repository.Like
}

func NewLikeService(likeRepo repository.Like) *LikeService {
	return &LikeService{likeRepo: likeRepo}
}

func (s *LikeService) LikePost(userId, postId int) error {
	count, err := s.likeRepo.GetUsersPostLikesCount(userId, postId)
	if err != nil {
		return err
	}
	if count > 0 {
		return UserLikePostErr{fmt.Sprintf("user %d already liked post %d", userId, postId)}
	}

	return s.likeRepo.LikePost(userId, postId)
}

func (s *LikeService) UnlikePost(userId, postId int) error {
	count, err := s.likeRepo.GetUsersPostLikesCount(userId, postId)
	if err != nil {
		return err
	}
	if count <= 0 {
		return UserLikePostErr{fmt.Sprintf("user %d didn't like post %d", userId, postId)}
	}
	return s.likeRepo.UnlikePost(userId, postId)
}
