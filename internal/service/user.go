package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"grates/internal/domain"
	"grates/internal/repository"
	"grates/pkg/auth"
	"math/rand"
	"time"
)

type UserService struct {
	repo         repository.User
	sigingKey    string
	passwordSalt string

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(repo repository.User, sigingKey, pswrdSalt string, accessTokenTTL, refreshTokenTTL time.Duration) *UserService {
	return &UserService{
		repo:            repo,
		sigingKey:       sigingKey,
		passwordSalt:    pswrdSalt,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

// CreateUser генерирует хэш пароля, сохраняет пользователя в БД.
// Возвращает int id созданного пользователя и ошибку.
func (s *UserService) CreateUser(user domain.UserSignUpInput) (int, error) {
	user.Password = auth.GeneratePasswordHash(user.Password, s.passwordSalt)
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUserById(id int) (domain.User, error) {
	return s.repo.GetUserById(id)
}

// Tokens структура по типу Double. Хранит пару access и refresh token
type Tokens struct {
	Access  string
	Refresh string
}

// AuthenticateUser получает пользователя из БД по заданным параметрам,
// возвращает сгенерированную пару токенов Tokens/
func (s *UserService) AuthenticateUser(email, password string) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)
	user, err := s.repo.GetUser(email, auth.GeneratePasswordHash(password, s.passwordSalt))
	logrus.Info(user)
	if err != nil {
		return tokens, err
	}

	return s.GenerateTokens(user)
}

// GenerateTokens , получая в качестве параметра domain.User, создает access и
// refresh токены, записывает соответствующий пользователю refresh token в БД.
// Возвращает пару access и refresh токеном Tokens.
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

// RefreshTokens ищет id пользователя по refresh токену, находит самого пользователя,
// возвращает сгенерированную пару access и refresh токенов Tokens.
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

// GetUserByEmail возвращает пользователя domain.User по уникальному email.
func (s *UserService) GetUserByEmail(email string) (domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

// GetAllUsers возвращет всех пользователей []domain.User
func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAllUsers()
}

// tokenClaims кастомный claims для access токена.
type tokenClaims struct {
	//User domain.User `json:"user"`
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

// ParseToken распаршивает access token и возвращает пользователя из claims'ов токена.
func (s *UserService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid siging method")
		}
		return []byte(s.sigingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, err
}

// newAccessToken генерирует новый access токен.
func (s *UserService) newAccessToken(user domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		tokenClaims{
			user.Id,
			//user,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			},
		},
	)
	return token.SignedString([]byte(s.sigingKey))
}

// newRefreshToken генерирует новый refresh токен.
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

// UpdateProfile обновляет профиль пользователя.
func (s *UserService) UpdateProfile(userId int, newProfile domain.ProfileUpdateInput) error {
	return s.repo.UpdateProfile(userId, newProfile)
}
