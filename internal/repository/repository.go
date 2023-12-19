package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"grates/internal/domain"
)

const (
	UsersTable      = "users"
	PostsTable      = "posts"
	CommentsTable   = "comments"
	LikesPostsTable = "likes_posts"
	AuthEmailsTable = "auth_emails"
	FriendsTable    = "friend_requests"
)

type User interface {
	CreateUser(user domain.UserSignUpInput) (int, error)
	GetUser(email string, password string) (domain.User, error)
	GetUserById(id int) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetAllUsers() ([]domain.User, error)

	SaveRefreshToken(userId int, session domain.Session) error
	GetUserIdByToken(refreshToken string) (int, error)

	UpdateProfile(userId int, newProfile domain.ProfileUpdateInput) error
}

type Post interface {
	Create(post domain.Post) (int, error)
	Get(postId int) (domain.Post, error)
	GetUsersPosts(userId int) ([]domain.Post, error)
	Update(id int, newPost domain.PostUpdateInput) error
	Delete(id int) error
}

type Comment interface {
	Create(comment domain.CommentCreateInput) (int, error)
	GetPostComments(postId int) ([]domain.Comment, error)
	Update(userId, commentId int, newComment domain.CommentUpdateInput) error
	Delete(userId, commentId int) error
}

type Like interface {
	GetPostLikesCount(postId int) (int, error)
	GetUsersPostLikesCount(userId, postId int) (int, error)
	LikePost(userId, postId int) error
	UnlikePost(userId, postId int) error
}

// TODO: fix
type Email interface {
	ReplaceEmail(userId int, hash string) error
	ConfirmEmail(userId int, hash string) error
}

type Friend interface {
	Get(id1, id2 int) (domain.Friend, error)
	GetFriendUsers(userId int) ([]domain.User, error)

	// GetRequests(userId int) ([]domain.FriendRequest, error)

	FriendRequest(fromId, toId int) error
	AcceptFriendRequest(fromId, toId int) error
	Unfriend(userId, friendId int) error
	Decline(userId, friendId int) error
}

type Repository struct {
	User    *UserRepository
	Post    *PostRepository
	Comment *CommentRepository
	Like    *LikeRepository
	Email   *EmailRepository
	Friend  *FriendRepository
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
		Like:    NewLikeRepository(db),
		Email:   NewEmailRepository(db),
		Friend:  NewFriendRepository(db),
	}
}
