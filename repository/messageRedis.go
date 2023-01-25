package repository

import (
	"context"

	"github.com/ponyo877/folks/entity"
)

type UserMySQLPresenter struct {
	ID          string
	DisplayName string
}

type UserMySQLPresenterList []UserMySQLPresenter

// userEntity
func userEntity(user *entity.User) *UserMySQLPresenter {
	return &UserMySQLPresenter{
		ID:          user.ID.String(),
		DisplayName: user.DisplayName.String(),
	}
}

// Append
func (r *MessageRepository) Append(roomID entity.UID, message *entity.Message) error {
	encodedMessage, err := entity.EncodeMessage(message)
	if err != nil {
		return err
	}
	return r.kvs.RPush(context.Background(), roomID.String(), encodedMessage).Err()
}

// ListRecent
func (r *MessageRepository) ListRecent(roomID entity.UID, size int64) ([]*entity.Message, error) {
	preDecodeMessages, err := r.kvs.LRange(context.Background(), roomID.String(), -size, -1).Result()
	if err != nil {
		return nil, err
	}
	var messageList []*entity.Message
	for _, preDecodeMessage := range preDecodeMessages {
		message, err := entity.DecodeMessage([]byte(preDecodeMessage))
		if err != nil {
			return nil, err
		}
		messageList = append(messageList, message)
	}
	return messageList, nil
}

// AddUser
func (r *MessageRepository) AddUser(session *entity.Session) error {
	roomID := session.RoomID.String()
	userID := session.User.ID.String()
	displayName := session.User.DisplayName.String()
	if _, err := r.kvs.HSet(context.Background(), "members_"+roomID, userID, displayName).Result(); err != nil {
		return err
	}
	return nil
}

// RemoveUser
func (r *MessageRepository) RemoveUser(session *entity.Session) error {
	roomID := session.RoomID.String()
	userID := session.User.ID.String()
	if _, err := r.kvs.HDel(context.Background(), "members_"+roomID, userID).Result(); err != nil {
		return err
	}
	return nil
}

// ListUser
func (r *MessageRepository) ListUser(roomID entity.UID) ([]*entity.User, error) {
	userMap, err := r.kvs.HGetAll(context.Background(), "members_"+roomID.String()).Result()
	if err != nil {
		return nil, err
	}
	var userList []*entity.User
	for userIDStr, displayNameStr := range userMap {
		userID, err := entity.StringToID(userIDStr)
		if err != nil {
			return nil, err
		}
		user := &entity.User{
			ID:          userID,
			DisplayName: entity.NewDisplayName(displayNameStr),
		}
		userList = append(userList, user)
	}
	return userList, nil
}
