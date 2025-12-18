package service

import (
	"context"
	"fmt"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/chat"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/chatmember"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/user"
)

type ChatService struct {
	client *ent.Client
}

func NewChatService(client *ent.Client) *ChatService {
	return &ChatService{client: client}
}

func (s *ChatService) CreateChat(ctx context.Context, name string, isGroup bool, creatorID int, memberIDs []int) (*ent.Chat, error) {
	// Create chat
	newChat, err := s.client.Chat.Create().
		SetName(name).
		SetIsGroup(isGroup).
		SetCreatorID(creatorID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	// Add creator as admin
	_, err = s.client.ChatMember.Create().
		SetChatID(newChat.ID).
		SetUserID(creatorID).
		SetIsAdmin(true).
		Save(ctx)
	if err != nil {
		// Rollback: delete the chat
		_ = s.client.Chat.DeleteOneID(newChat.ID).Exec(ctx)
		return nil, fmt.Errorf("failed to add creator as member: %w", err)
	}

	// Add other members
	for _, memberID := range memberIDs {
		if memberID == creatorID {
			continue // Skip creator as already added
		}

		// Check if user exists
		exists, _ := s.client.User.Query().
			Where(user.ID(memberID)).
			Exist(ctx)
		if !exists {
			continue // Skip non-existent users
		}

		_, err = s.client.ChatMember.Create().
			SetChatID(newChat.ID).
			SetUserID(memberID).
			SetIsAdmin(false).
			Save(ctx)
		if err != nil {
			// Continue even if adding a member fails
			continue
		}
	}

	return newChat, nil
}

func (s *ChatService) ListUserChats(ctx context.Context, userID, limit, offset int, withMembers bool) ([]*ent.Chat, error) {
	q := s.client.Chat.Query().
		Where(chat.HasMembersWith(chatmember.HasUserWith(user.ID(userID)))).
		WithCreator()

	if withMembers {
		q = q.WithMembers(func(q *ent.ChatMemberQuery) {
			q.WithUser()
		})
	}

	chats, err := q.Limit(limit).
		Offset(offset).
		Order(ent.Desc(chat.FieldUpdatedAt)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list chats: %w", err)
	}

	return chats, nil
}

func (s *ChatService) GetChatByID(ctx context.Context, chatID int) (*ent.Chat, error) {
	chatEntity, err := s.client.Chat.Query().
		Where(chat.ID(chatID)).
		WithCreator().
		WithMembers(func(q *ent.ChatMemberQuery) {
			q.WithUser()
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return chatEntity, nil
}

func (s *ChatService) IsUserMemberOfChat(ctx context.Context, chatID, userID int) (bool, error) {
	isMember, err := s.client.ChatMember.Query().
		Where(
			chatmember.HasChatWith(chat.ID(chatID)),
			chatmember.HasUserWith(user.ID(userID)),
		).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}

	return isMember, nil
}

func (s *ChatService) IsUserAdminOfChat(ctx context.Context, chatID, userID int) (bool, error) {
	member, err := s.client.ChatMember.Query().
		Where(
			chatmember.HasChatWith(chat.ID(chatID)),
			chatmember.HasUserWith(user.ID(userID)),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	return member.IsAdmin, nil
}

func (s *ChatService) UpdateChatName(ctx context.Context, chatID int, name string) (*ent.Chat, error) {
	chatEntity, err := s.client.Chat.UpdateOneID(chatID).
		SetName(name).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to update chat: %w", err)
	}

	return chatEntity, nil
}

func (s *ChatService) DeleteChat(ctx context.Context, chatID int) error {
	err := s.client.Chat.DeleteOneID(chatID).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("chat not found")
		}
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (s *ChatService) GetChatCreator(ctx context.Context, chatID int) (*ent.User, error) {
	chatEntity, err := s.client.Chat.Get(ctx, chatID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("chat not found")
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	creator, err := chatEntity.QueryCreator().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("creator not found")
		}
		return nil, fmt.Errorf("failed to get creator: %w", err)
	}

	return creator, nil
}

func (s *ChatService) AddMembers(ctx context.Context, chatID int, memberIDs []int) error {
	for _, memberID := range memberIDs {
		// Check if already a member
		exists, _ := s.client.ChatMember.Query().
			Where(
				chatmember.HasChatWith(chat.ID(chatID)),
				chatmember.HasUserWith(user.ID(memberID)),
			).
			Exist(ctx)
		if exists {
			continue // Skip if already a member
		}

		// Check if user exists
		userExists, _ := s.client.User.Query().
			Where(user.ID(memberID)).
			Exist(ctx)
		if !userExists {
			continue // Skip non-existent users
		}

		// Add member
		_, err := s.client.ChatMember.Create().
			SetChatID(chatID).
			SetUserID(memberID).
			SetIsAdmin(false).
			Save(ctx)
		if err != nil {
			// Continue even if adding fails
			continue
		}
	}

	return nil
}

func (s *ChatService) RemoveMember(ctx context.Context, chatID, memberID int) error {
	// Delete member
	_, err := s.client.ChatMember.Delete().
		Where(
			chatmember.HasChatWith(chat.ID(chatID)),
			chatmember.HasUserWith(user.ID(memberID)),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}
