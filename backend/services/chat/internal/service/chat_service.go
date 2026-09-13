package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/model"
	"github.com/williamsebastianliman/WEB-WS-242/services/chat/internal/repository"
)

type ChatService interface {
	CreateChat(groupID, senderID uuid.UUID, message string) error
	CreateGroup(memberIDs []uuid.UUID) (uuid.UUID, error)
	AddMember(groupID, userID uuid.UUID) error
	CreateConversation(userA, userB uuid.UUID) (uuid.UUID, error)
	GetAllGroupsByUser(userID uuid.UUID) ([]model.Group, error)
	GetAllChatsByGroup(groupID uuid.UUID) ([]model.Chat, error)
	DeleteChat(chatID uuid.UUID) error
	MaskChat(chatID uuid.UUID) error
	RemoveMember(groupID, userID uuid.UUID) error
	GetAllMembersByGroup(groupID uuid.UUID) ([]model.GroupMember, error)
}

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}

func (s *chatService) CreateChat(groupID, senderID uuid.UUID, message string) error {
	chat := &model.Chat{
		GroupID:   groupID,
		SenderID:  senderID,
		Message:   message,
		IsMasked:  false,
		CreatedAt: time.Now(),
	}
	return s.repo.CreateChat(chat)
}

func (s *chatService) CreateGroup(memberIDs []uuid.UUID) (uuid.UUID, error) {
	group := &model.Group{
		IsPrivate: false,
	}
	err := s.repo.CreateGroup(group)
	if err != nil {
		return uuid.Nil, err
	}
	for _, memberID := range memberIDs {
		if err := s.repo.AddMember(group.GroupID, memberID); err != nil {
			return uuid.Nil, err
		}
	}
	return group.GroupID, nil
}

func (s *chatService) AddMember(groupID, userID uuid.UUID) error {
	return s.repo.AddMember(groupID, userID)
}

func (s *chatService) CreateConversation(userA, userB uuid.UUID) (uuid.UUID, error) {
	return s.repo.CreateConversation(userA, userB)
}

func (s *chatService) GetAllGroupsByUser(userID uuid.UUID) ([]model.Group, error) {
	return s.repo.GetAllGroupsByUser(userID)
}

func (s *chatService) GetAllChatsByGroup(groupID uuid.UUID) ([]model.Chat, error) {
	return s.repo.GetAllChatsByGroup(groupID)
}

func (s *chatService) DeleteChat(chatID uuid.UUID) error {
	return s.repo.DeleteChat(chatID)
}

func (s *chatService) MaskChat(chatID uuid.UUID) error {
	return s.repo.MaskChat(chatID)
}

func (s *chatService) RemoveMember(groupID, userID uuid.UUID) error {
	return s.repo.RemoveMember(groupID, userID)
}

func (s *chatService) GetAllMembersByGroup(groupID uuid.UUID) ([]model.GroupMember, error) {
    return s.repo.GetAllMembersByGroup(groupID)
}