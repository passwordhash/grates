package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"grates/internal/repository"
)

type FriendService struct {
	friendRepo repository.Friend
}

func NewFriendService(friendRepo repository.Friend) *FriendService {
	return &FriendService{friendRepo: friendRepo}
}

func (f *FriendService) SendFriendRequest(fromId, toId int) error {
	if err := f.checkIds(fromId, toId); err != nil {
		return err
	}

	return f.friendRepo.FriendRequest(fromId, toId)
}

func (f *FriendService) AcceptFriendRequest(fromId, toId int) error {
	if err := f.checkIds(fromId, toId); err != nil {
		return err
	}

	request, err := f.friendRepo.Get(fromId, toId)
	if err != nil {
		return fmt.Errorf("can't get friend request: %w", err)
	}

	if request.IsConfirmed {
		return errors.New("you can't accept friend request, because it's already confirmed")
	}

	if request.FromId != fromId {
		return errors.New("you can't accept friend request, because you are not recipient")
	}

	return f.friendRepo.AcceptFriendRequest(fromId, toId)
}

func (f *FriendService) Unfriend(userId, friendId int) error {
	if err := f.checkIds(userId, friendId); err != nil {
		return err
	}

	request, err := f.friendRepo.Get(userId, friendId)
	if err != nil {
		return fmt.Errorf("can't get friend request: %w", err)
	}

	// Если заявка подтверждена (т. е. пользователи являются друзьями),
	// то удаляем запись из таблицы и добавляем обратно, но с другими id
	// (как будто пользователь, которого удалили, отправил заявку)
	if request.IsConfirmed {
		return f.friendRepo.Unfriend(userId, friendId)
	}

	// Если заявка не подтверждена и отправитель заявки удаляет её, удалить запись из таблицы.
	// Тип "Я отправил запрос, но потом передумал"
	if request.FromId == userId {
		return f.friendRepo.Decline(userId, friendId)
	}

	// Если заявка не подтверждена и получатель заявки хочет удалить,
	// в данной версии приложения ничего не будет проихсодить.
	// TODO: функционал подписок
	if request.ToId == userId {
		// TEMP
		logrus.Infof("user %d wants to delete friend request from user %d, but nothing to happend", userId, friendId)
	}

	return nil
}

func (f *FriendService) checkIds(id1, id2 int) error {
	if id1 == id2 {
		return errors.New("you can't send friend request to yourself")
	}

	return nil
}
