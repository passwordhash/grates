package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"grates/internal/domain"
	"grates/internal/repository"
	"math/rand"
	"time"
)

const (
	// TODO: remove!
	salt = "hjqrhjqw124617ajfhajs"
)

type UserService struct {
	repo      repository.User
	sigingKey string

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repo repository.User, sigingKey string, accessTokenTTL, refreshTokenTTL time.Duration) *UserService {
	return &UserService{
		repo:            repo,
		sigingKey:       sigingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UserService) CreateUser(user domain.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

type Tokens struct {
	Access  string
	Refresh string
}

func (s *UserService) AuthenticateUser(email, password string) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return tokens, err
	}

	return s.GenerateTokens(user)
}

func (s *UserService) GenerateTokens(user domain.User) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens.Access, err = s.newAccessToken(user)
	if err != nil {
		return tokens, err
	}

	tokens.Refresh, err = s.newRefreshToken()
	if err != nil {
		return tokens, nil
	}

	err = s.repo.SaveRefreshToken(user.Id, domain.Session{
		RefreshToken: tokens.Refresh,
		//ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
		TTL: s.refreshTokenTTL,
	})
	if err != nil {
		return Tokens{}, err
	}

	return tokens, err
}

func (s *UserService) RefreshTokens(refreshToken string) (Tokens, error) {
	userId, err := s.repo.GetUserIdByToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetUserById(userId)
	if err != nil {
		return Tokens{}, err
	}

	return s.GenerateTokens(user)
}

func (s *UserService) GetUserByEmail(email string) (domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

type tokenClaims struct {
	User domain.User `json:"user"`
	jwt.RegisteredClaims
}

func (s *UserService) ParseToken(accessToken string) (domain.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid siging method")
		}
		return []byte(s.sigingKey), nil
	})

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return domain.User{}, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.User, err
}

func (s *UserService) newAccessToken(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		tokenClaims{
			user,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			},
		},
	)
	return token.SignedString([]byte(s.sigingKey))
}

func (s *UserService) newRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
