package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"grates/internal/domain"
)

type User interface {
	CreateUser(user domain.User) (int, error)
	GetUser(email string, password string) (domain.User, error)
	GetUserById(id int) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)

	SaveRefreshToken(userId int, session domain.Session) error
	GetUserIdByToken(refreshToken string) (int, error)
}

type Post interface {
	Create(post domain.Post) (int, error)
	Get(postId int) (domain.Post, error)
	GetUsersPosts(postId int) ([]domain.Post, error)
	Update(id int, newPost domain.PostUpdateInput) error
	Delete(id int) error
}

type Comment interface {
	Create(comment domain.CommentCreateInput) (int, error)
	GetPostComments(postId int) ([]domain.Comment, error)
	Update(id int, newComment domain.CommentCreateInput) error
	Delete(id int) error
}

type Repository struct {
	User    *UserRepository
	Post    *PostRepository
	Comment *CommentRepository
}

type UserRepository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewUserRepository(db *sqlx.DB, rdb *redis.Client) *UserRepository {
	return &UserRepository{db: db, rdb: rdb}
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		User:    NewUserRepository(db, rdb),
		Post:    NewPostPostgres(db),
		Comment: NewCommentRepository(db),
	}
}
