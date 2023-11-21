package repository

import (
	"context"
	"github.com/sirupsen/logrus"
	"grates/internal/domain"
	"strconv"
)

func (r *UserRepository) SetUser(user domain.User) {
	_, _ = r.rdb.HSet(context.Background(), strconv.Itoa(user.Id), "email", user.Email).Result()
	v := r.rdb.HGetAll(context.Background(), strconv.Itoa(user.Id))
	logrus.Info(v)
}

func (r *UserRepository) SaveRefreshToken(userId int, session domain.Session) error {
	c := context.Background()
	logrus.Info(userId)
	_, err := r.rdb.Set(c, session.RefreshToken, userId, session.TTL).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserIdByToken(refreshToken string) (int, error) {
	var id int
	c := context.Background()

	v, err := r.rdb.Get(c, refreshToken).Result()
	if err != nil {
		return 0, err
	}

	id, err = strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return id, nil
}
