package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
)

type User interface {
	CreateUser(user domain.User) (int, error)
	GetUserByEmail(email string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	AuthenticateUser(email string, password string) (string, error)
	ParseToken(token string) (domain.User, error)
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
