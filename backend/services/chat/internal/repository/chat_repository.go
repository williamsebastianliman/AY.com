package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateChat(chat *model.Chat) error
	CreateGroup(group *model.Group) error
	AddMember(groupID, userID uuid.UUID) error
	CreateConversation(userA, userB uuid.UUID) (uuid.UUID, error)
	GetAllGroupsByUser(userID uuid.UUID) ([]model.Group, error)
	GetAllChatsByGroup(groupID uuid.UUID) ([]model.Chat, error)
	DeleteChat(chatID uuid.UUID) error
	MaskChat(chatID uuid.UUID) error
	RemoveMember(groupID, userID uuid.UUID) error
	GetAllMembersByGroup(groupID uuid.UUID) ([]model.GroupMember, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateChat(chat *model.Chat) error {
	chat.ChatID = uuid.New()
	chat.CreatedAt = time.Now()
	return r.db.Create(chat).Error
}

func (r *chatRepository) CreateGroup(group *model.Group) error {
	group.GroupID = uuid.New()
	return r.db.Create(&group).Error
}

func (r *chatRepository) AddMember(groupID, userID uuid.UUID) error {
	member := model.GroupMember{GroupID: groupID, MemberID: userID}
	return r.db.Create(&member).Error
}

func (r *chatRepository) CreateConversation(userA, userB uuid.UUID) (uuid.UUID, error) {
	group := &model.Group{IsPrivate: true}
	if err := r.CreateGroup(group); err != nil {
		return uuid.Nil, err
	}
	if err := r.AddMember(group.GroupID, userA); err != nil {
		return uuid.Nil, err
	}
	if err := r.AddMember(group.GroupID, userB); err != nil {
		return uuid.Nil, err
	}
	return group.GroupID, nil
}

func (r *chatRepository) GetAllGroupsByUser(userID uuid.UUID) ([]model.Group, error) {
	var groups []model.Group
	err := r.db.Table("groups").
		Joins("JOIN group_members ON groups.group_id = group_members.group_id").
		Where("group_members.member_id = ?", userID).
		Find(&groups).Error
	return groups, err
}

func (r *chatRepository) GetAllChatsByGroup(groupID uuid.UUID) ([]model.Chat, error) {
	var chats []model.Chat
	err := r.db.Where("group_id = ?", groupID).Order("created_at ASC").Find(&chats).Error
	return chats, err
}

func (r *chatRepository) DeleteChat(chatID uuid.UUID) error {
	return r.db.Delete(&model.Chat{}, "chat_id = ?", chatID).Error
}

func (r *chatRepository) MaskChat(chatID uuid.UUID) error {
	return r.db.Model(&model.Chat{}).Where("chat_id = ?", chatID).Update("is_masked", true).Error
}

func (r *chatRepository) RemoveMember(groupID, userID uuid.UUID) error {
	return r.db.Delete(&model.GroupMember{}, "group_id = ? AND member_id = ?", groupID, userID).Error
}

func (r *chatRepository) GetAllMembersByGroup(groupID uuid.UUID) ([]model.GroupMember, error) {
    var members []model.GroupMember
    err := r.db.Where("group_id = ?", groupID).Find(&members).Error
    return members, err
}