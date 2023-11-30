package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
	"time"
)

type User interface {
	CreateUser(user domain.User) (int, error)
	GetUserByEmail(email string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	AuthenticateUser(email string, password string) (Tokens, error)
	GenerateTokens(user domain.User) (Tokens, error)
	ParseToken(token string) (domain.User, error)
	RefreshTokens(refreshToken string) (Tokens, error)
}

type Post interface {
	Create(post domain.Post) (int, error)
	Get(postId int) (domain.Post, error)
	GetUsersPosts(userId int) ([]domain.Post, error)
	Update(id int, newPost domain.PostUpdateInput) error
	Delete(id int) error
}

type Comment interface {
}

type Service struct {
	User
	Post
}

type Deps struct {
	SigingKey    string
	PasswordSalt string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewService(repos *repository.Repository, deps Deps) *Service {
	return &Service{
		User: NewUserService(repos.User, deps.SigingKey, deps.PasswordSalt, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Post: NewPostService(repos.Post),
	}
}
