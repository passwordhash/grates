package service

import "grates/internal/repository"

type LikeService struct {
	likeRepo repository.Like
}

func NewLikeService(likeRepo repository.Like) *LikeService {
	return &LikeService{likeRepo: likeRepo}
}

func (s *LikeService) LikePost(userId, postId int) error {
	return s.likeRepo.LikePost(userId, postId)
}

func (s *LikeService) UnlikePost(userId, postId int) error {
	return s.likeRepo.UnlikePost(userId, postId)
}
