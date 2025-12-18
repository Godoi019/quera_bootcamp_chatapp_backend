package service

import (
	"context"
	"fmt"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/chat"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/message"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/user"
)

type MessageService struct {
	client *ent.Client
}

func NewMessageService(client *ent.Client) *MessageService {
	return &MessageService{client: client}
}

func (s *MessageService) SendMessage(ctx context.Context, chatID, senderID int, content string) (*ent.Message, error) {
	// Create message
	newMessage, err := s.client.Message.Create().
		SetChatID(chatID).
		SetSenderID(senderID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return newMessage, nil
}

func (s *MessageService) GetMessageByID(ctx context.Context, messageID int) (*ent.Message, error) {
	msg, err := s.client.Message.Query().
		Where(message.ID(messageID)).
		WithSender().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return msg, nil
}

func (s *MessageService) ListChatMessages(ctx context.Context, chatID, limit, offset int) ([]*ent.Message, error) {
	messages, err := s.client.Message.Query().
		Where(message.HasChatWith(chat.ID(chatID))).
		WithSender().
		Order(ent.Desc(message.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return messages, nil
}

func (s *MessageService) UpdateMessage(ctx context.Context, messageID int, content string) (*ent.Message, error) {
	updatedMessage, err := s.client.Message.UpdateOneID(messageID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	return updatedMessage, nil
}

func (s *MessageService) DeleteMessage(ctx context.Context, messageID int) error {
	err := s.client.Message.DeleteOneID(messageID).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("message not found")
		}
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *MessageService) GetMessageSender(ctx context.Context, messageID int) (*ent.User, error) {
	msg, err := s.client.Message.Get(ctx, messageID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	sender, err := msg.QuerySender().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("sender not found")
		}
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	return sender, nil
}

func (s *MessageService) GetMessageChat(ctx context.Context, messageID int) (*ent.Chat, error) {
	msg, err := s.client.Message.Get(ctx, messageID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	chatEntity, err := msg.QueryChat().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return chatEntity, nil
}

func (s *MessageService) IsUserSenderOfMessage(ctx context.Context, messageID, userID int) (bool, error) {
	exists, err := s.client.Message.Query().
		Where(
			message.ID(messageID),
			message.HasSenderWith(user.ID(userID)),
		).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check sender: %w", err)
	}

	return exists, nil
}
