package service

import (
	"grates/internal/entity"
	"grates/pkg/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
}

type User interface {
	GetAllUsers() ([]entity.User, error)
}

type Post interface {
}

type Comment interface {
}

type Service struct {
	Authorization
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
	}
}
