package service

import (
	"errors"
	"github.com/cha1l/sayrsa-2.0/models"
	"github.com/cha1l/sayrsa-2.0/pkg/repository"
	"time"
)

const (
	UpdateTokenTime = 24 * time.Hour
)

type ConversationService struct {
	repo repository.Conversations
}

func NewConversationService(repo repository.Conversations) *ConversationService {
	return &ConversationService{
		repo: repo,
	}
}

func (s *ConversationService) GetConversationInfo(username string, convID int) (*models.Conversation, error) {
	info, err := s.repo.GetConversationInfo(convID)
	if err != nil {
		return nil, err
	}
	return info, s.UpdateToken(username)
}

func (s *ConversationService) CreateConversation(username string, title string, members []string) (int, []models.PublicKey, error) {
	convID, err := s.repo.CreateConversation(title, members)
	if err != nil {
		return 0, nil, err
	}
	publicKeys, err := s.repo.GetUsersPublicKeys(members...)
	if err != nil {
		return 0, nil, err
	}

	return convID, publicKeys, s.UpdateToken(username)
}

func (s *ConversationService) UpdateToken(username string) error {
	token, err := s.repo.GetUserToken(username)
	if err != nil {
		return err
	}

	token.ExpiresAt = token.ExpiresAt.Add(UpdateTokenTime)

	return s.repo.UpdateUserToken(token)
}

func (s *ConversationService) GetPublicKey(username string) (string, error) {
	publicKeys, err := s.repo.GetUsersPublicKeys(username)
	if len(publicKeys) == 1 {
		publicKey := publicKeys[0]
		return publicKey.PublicKey, err
	}
	return "", errors.New("wrong length of slice")
}
