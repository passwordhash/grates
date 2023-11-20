package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"grates/internal/entity"
	"grates/pkg/repository"
	"os"
	"time"
)

const (
	salt     = "hjqrhjqw124617ajfhajs"
	tokenTTL = 12 * time.Hour
)

type UserService struct {
	repo      repository.User
	sigingKey string
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo, sigingKey: os.Getenv("JWT_SIGING_KEY")}
}

func (s *UserService) CreateUser(user entity.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *UserService) AuthenticateUser(email string, password string) (string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	return s.generateToken(user)
}

func (s *UserService) GetUserByEmail(email string) (entity.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetAllUsers() ([]entity.User, error) {
	return s.repo.GetAllUsers()
}

type tokenClaims struct {
	User entity.User `json:"user"`
	jwt.RegisteredClaims
}

func (s *UserService) ParseToken(accessToken string) (entity.User, error) {
	var user entity.User

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid siging method")
		}
		return []byte(s.sigingKey), nil
	})

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return user, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.User, err
}

func (s *UserService) generateToken(user entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		tokenClaims{
			user,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			},
		},
	)
	return token.SignedString([]byte(s.sigingKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
