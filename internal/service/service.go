package service

import (
	"grates/internal/domain"
	"grates/internal/repository"
	"time"
)

type User interface {
	CreateUser(user domain.UserSignUpInput) (int, error)
	GetUserById(id int) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)

	AuthenticateUser(email string, password string) (Tokens, error)
	GenerateTokens(user domain.User) (Tokens, error)
	ParseToken(token string) (int, error)
	RefreshTokens(refreshToken string) (Tokens, error)

	UpdateProfile(userId int, newProfile domain.ProfileUpdateInput) error
}

type Post interface {
	Create(post domain.Post) (int, error)
	GetWithAdditions(postId int) (domain.Post, error)
	GetUsersPosts(userId int) ([]domain.Post, error)
	Update(id int, newPost domain.PostUpdateInput) error
	Delete(id int) error
	IsPostBelongsToUser(userId, postId int) (bool, error)
}

type Comment interface {
	Create(comment domain.CommentCreateInput) (int, error)
	GetPostComments(postId int) ([]domain.Comment, error)
	Delete(userId, commentId int) error
	Update(userId, commentId int, newComment domain.CommentUpdateInput) error
}

type Email interface {
	ReplaceConfirmationEmail(userId int, to, name string) error
	ConfirmEmail(hash string) error
	sendAuthEmail(to, name, hash string) error
}

type Like interface {
	LikePost(userId, postId int) error
	UnlikePost(userId, postId int) error
}

type Friend interface {
	GetFriends(userId int) ([]domain.User, error)
	SendFriendRequest(fromId, toId int) error
	AcceptFriendRequest(fromId, toId int) error
	Unfriend(userId, friendId int) error
}

type Service struct {
	User
	Post
	Comment
	Like
	Email
	Friend
}

type Deps struct {
	SigingKey    string
	PasswordSalt string

	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	EmailDeps
}

type EmailDeps struct {
	SmtpHost string
	SmtpPort int

	From     string
	Password string

	BaseUrl string
}

func NewService(repos *repository.Repository, deps Deps) *Service {
	return &Service{
		User:    NewUserService(repos.User, deps.SigingKey, deps.PasswordSalt, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Post:    NewPostService(repos.Post, repos.Comment, repos.Like),
		Comment: NewCommentService(repos.Comment),
		Like:    NewLikeService(repos.Like),
		// TODO: fix (pointer)
		Email:  NewEmailService(*repos.Email, deps.EmailDeps),
		Friend: NewFriendService(repos.Friend),
	}
}
