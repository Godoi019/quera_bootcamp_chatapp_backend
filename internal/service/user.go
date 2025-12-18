package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/auth"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent/user"
)

type UserService struct {
	client      *ent.Client
	authService *auth.Service
}

func NewUserService(client *ent.Client, authService *auth.Service) *UserService {
	return &UserService{
		client:      client,
		authService: authService,
	}
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*ent.User, error) {
	query := s.client.User.Query().
		Limit(limit).
		Offset(offset).
		Order(ent.Asc(user.FieldUsername))

	users, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*ent.User, error) {
	u, err := s.client.User.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return u, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*ent.User, error) {
	u, err := s.client.User.Query().
		Where(user.Username(username)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return u, nil
}

func (s *UserService) UserExists(ctx context.Context, username string) (bool, error) {
	exists, err := s.client.User.Query().
		Where(user.Username(username)).
		Exist(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

func (s *UserService) CreateUser(ctx context.Context, username, password, displayName string) (*ent.User, error) {
	// Hash password
	hashedPassword, err := s.authService.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	newUser, err := s.client.User.Create().
		SetUsername(username).
		SetPassword(hashedPassword).
		SetDisplayName(displayName).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, displayName, password string) (*ent.User, error) {
	update := s.client.User.UpdateOneID(id)

	if displayName != "" {
		update.SetDisplayName(displayName)
	}

	if password != "" {
		hashedPassword, err := s.authService.HashPassword(password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		update.SetPassword(hashedPassword)
	}

	u, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.client.User.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) UpdateLastSeen(ctx context.Context, userID int) error {
	now := time.Now()
	err := s.client.User.UpdateOneID(userID).
		SetLastSeen(now).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	return nil
}

func (s *UserService) VerifyPassword(hashedPassword, password string) error {
	return s.authService.VerifyPassword(hashedPassword, password)
}

func (s *UserService) CreateToken(userID int, username string) (string, error) {
	token, err := s.authService.CreateToken(userID, username)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return token, nil
}
