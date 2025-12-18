package handler

import (
	"context"
	"time"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/auth"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/model"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/service"
	f "github.com/Hossara/quera_bootcamp_chatapp_backend/pkg/fiber"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(client *ent.Client, authService *auth.Service) *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(client, authService),
	}
}

func (h *UserHandler) ListUsers(c fiber.Ctx) error {
	limit := utils.QueryInt(c, "limit", 50)
	offset := utils.QueryInt(c, "offset", 0)

	users, err := h.userService.ListUsers(context.Background(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to list users",
		})
	}

	userProfiles := make([]model.UserProfile, len(users))
	for i, u := range users {
		userProfiles[i] = model.UserProfile{
			ID:          u.ID,
			Username:    u.Username,
			DisplayName: u.DisplayName,
			CreatedAt:   u.CreatedAt,
			LastSeen:    u.LastSeen,
		}
	}

	return c.JSON(userProfiles)
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	id, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid user id",
		})
	}

	u, err := h.userService.GetUserByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error: "user not found",
		})
	}

	return c.JSON(model.UserProfile{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		CreatedAt:   u.CreatedAt,
		LastSeen:    u.LastSeen,
	})
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid user id",
		})
	}

	if userID != id {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "you can only update your own profile",
		})
	}

	req := new(model.UpdateUserRequest)
	if err := f.ParseRequestBody(c, req); err != nil {
		return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
	}

	if req.Password != "" && len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "password must be at least 6 characters",
		})
	}

	u, err := h.userService.UpdateUser(context.Background(), id, req.DisplayName, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to update user",
		})
	}

	return c.JSON(model.UserProfile{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		CreatedAt:   u.CreatedAt,
		LastSeen:    u.LastSeen,
	})
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	id, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid user id",
		})
	}

	if userID != id {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "you can only delete your own profile",
		})
	}

	err = h.userService.DeleteUser(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to delete user",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *UserHandler) UpdateLastSeen(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	err := h.userService.UpdateLastSeen(context.Background(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to update last seen",
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"last_seen": time.Now(),
	})
}
