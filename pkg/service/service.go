package service

import (
	"grates/internal/entity"
	"grates/pkg/repository"
)

type User interface {
	CreateUser(user entity.User) (int, error)
	GetUserByEmail(email string) (entity.User, error)
	GetAllUsers() ([]entity.User, error)
}

type Post interface {
}

type Comment interface {
}

type Service struct {
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repos.User),
	}
}
