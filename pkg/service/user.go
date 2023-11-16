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

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserByEmail(email string) (entity.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetAllUsers() ([]entity.User, error) {
	return s.repo.GetAllUsers()
}

func generatePasswordHash(p string) string {
	hash := sha1.New()
	hash.Write([]byte(p))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
