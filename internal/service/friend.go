package service

import (
	"errors"
	"grates/internal/repository"
)

type FriendService struct {
	friendRepo repository.Friend
}

func NewFriendService(friendRepo repository.Friend) *FriendService {
	return &FriendService{friendRepo: friendRepo}
}

func (f *FriendService) SendFriendRequest(fromId, toId int) error {
	if fromId == toId {
		return errors.New("you can't send friend request to yourself")
	}

	return f.friendRepo.FriendRequest(fromId, toId)
}

func (f *FriendService) AcceptFriendRequest(id1, id2 int) error {
	if id1 == id2 {
		return errors.New("you can't accept friend request from yourself")
	}

	return f.friendRepo.AcceptFriendRequest(id1, id2)
}
