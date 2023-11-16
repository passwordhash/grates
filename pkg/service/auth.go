package service

import (
	"crypto/sha1"
	"fmt"
	"grates/internal/entity"
	"grates/pkg/repository"
)

const (
	salt = "blk;jgklr5345j34fbn245j"
)

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(p string) string {
	hash := sha1.New()
	hash.Write([]byte(p))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
