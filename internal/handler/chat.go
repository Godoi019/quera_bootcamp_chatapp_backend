package handler

import (
	"context"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/model"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/repository/ent"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/internal/service"
	f "github.com/Hossara/quera_bootcamp_chatapp_backend/pkg/fiber"
	"github.com/Hossara/quera_bootcamp_chatapp_backend/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(client *ent.Client) *ChatHandler {
	return &ChatHandler{
		chatService: service.NewChatService(client),
	}
}

func (h *ChatHandler) CreateChat(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	req := new(model.CreateChatRequest)
	if err := f.ParseRequestBody(c, req); err != nil {
		return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
	}

	if len(req.MemberIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "at least one member is required",
		})
	}

	// Create chat
	newChat, err := h.chatService.CreateChat(context.Background(), req.Name, req.IsGroup, userID, req.MemberIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to create chat",
		})
	}

	creatorID := userID
	if newChat.Edges.Creator != nil {
		creatorID = newChat.Edges.Creator.ID
	}

	return c.Status(fiber.StatusCreated).JSON(model.ChatResponse{
		ID:        newChat.ID,
		Name:      newChat.Name,
		IsGroup:   newChat.IsGroup,
		CreatorID: creatorID,
		CreatedAt: newChat.CreatedAt,
		UpdatedAt: newChat.UpdatedAt,
	})
}

func (h *ChatHandler) ListChats(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	withMembers := fiber.Query[bool](c, "with_members", false)
	limit := utils.QueryInt(c, "limit", 50)
	offset := utils.QueryInt(c, "offset", 0)

	chats, err := h.chatService.ListUserChats(context.Background(), userID, limit, offset, withMembers)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to list chats",
		})
	}

	chatResponses := make([]model.ChatDetailResponse, 0, len(chats))
	for _, chat := range chats {
		members := make([]model.ChatMemberResponse, 0, len(chat.Edges.Members))
		for _, member := range chat.Edges.Members {
			if member.Edges.User != nil {
				members = append(members, model.ChatMemberResponse{
					UserID:   member.Edges.User.ID,
					Username: member.Edges.User.Username,
					IsAdmin:  member.IsAdmin,
					JoinedAt: member.JoinedAt,
				})
			}
		}

		creatorID := 0
		if chat.Edges.Creator != nil {
			creatorID = chat.Edges.Creator.ID
		}

		chatResponses = append(chatResponses, model.ChatDetailResponse{
			ChatResponse: model.ChatResponse{
				ID:        chat.ID,
				Name:      chat.Name,
				IsGroup:   chat.IsGroup,
				CreatorID: creatorID,
				CreatedAt: chat.CreatedAt,
				UpdatedAt: chat.UpdatedAt,
			},
			Members: members,
		})
	}

	return c.JSON(chatResponses)
}

func (h *ChatHandler) GetChat(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid chat id",
		})
	}

	// Check if user is a member
	isMember, err := h.chatService.IsUserMemberOfChat(context.Background(), chatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to check membership",
		})
	}
	if !isMember {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "you are not a member of this chat",
		})
	}

	chatEntity, err := h.chatService.GetChatByID(context.Background(), chatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error: "chat not found",
		})
	}

	members := make([]model.ChatMemberResponse, 0, len(chatEntity.Edges.Members))
	for _, member := range chatEntity.Edges.Members {
		if member.Edges.User != nil {
			members = append(members, model.ChatMemberResponse{
				UserID:   member.Edges.User.ID,
				Username: member.Edges.User.Username,
				IsAdmin:  member.IsAdmin,
				JoinedAt: member.JoinedAt,
			})
		}
	}

	creatorID := 0
	if chatEntity.Edges.Creator != nil {
		creatorID = chatEntity.Edges.Creator.ID
	}

	return c.JSON(model.ChatDetailResponse{
		ChatResponse: model.ChatResponse{
			ID:        chatEntity.ID,
			Name:      chatEntity.Name,
			IsGroup:   chatEntity.IsGroup,
			CreatorID: creatorID,
			CreatedAt: chatEntity.CreatedAt,
			UpdatedAt: chatEntity.UpdatedAt,
		},
		Members: members,
	})
}

func (h *ChatHandler) UpdateChat(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid chat id",
		})
	}

	// Check if user is admin
	isAdmin, err := h.chatService.IsUserAdminOfChat(context.Background(), chatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to check admin status",
		})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "only admins can update the chat",
		})
	}

	req := new(model.UpdateChatRequest)
	if err := f.ParseRequestBody(c, req); err != nil {
		return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
	}

	chatEntity, err := h.chatService.UpdateChatName(context.Background(), chatID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to update chat",
		})
	}

	creatorID := 0
	if chatEntity.Edges.Creator != nil {
		creatorID = chatEntity.Edges.Creator.ID
	}

	return c.JSON(model.ChatResponse{
		ID:        chatEntity.ID,
		Name:      chatEntity.Name,
		IsGroup:   chatEntity.IsGroup,
		CreatorID: creatorID,
		CreatedAt: chatEntity.CreatedAt,
		UpdatedAt: chatEntity.UpdatedAt,
	})
}

func (h *ChatHandler) DeleteChat(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid chat id",
		})
	}

	// Check if user is the creator
	creator, err := h.chatService.GetChatCreator(context.Background(), chatID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error: "chat not found",
		})
	}

	if creator.ID != userID {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "only the creator can delete the chat",
		})
	}

	err = h.chatService.DeleteChat(context.Background(), chatID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to delete chat",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *ChatHandler) AddMembers(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid chat id",
		})
	}

	// Check if user is admin
	isAdmin, err := h.chatService.IsUserAdminOfChat(context.Background(), chatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to check membership",
		})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "only admins can add members",
		})
	}

	req := new(model.AddMembersRequest)
	if err := f.ParseRequestBody(c, req); err != nil {
		return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
	}

	err = h.chatService.AddMembers(context.Background(), chatID, req.MemberIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to add members",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "members added successfully",
	})
}

func (h *ChatHandler) RemoveMember(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid chat id",
		})
	}

	memberID, err := utils.ParamsInt(c, "memberId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid member id",
		})
	}

	// Check if user is admin
	isAdmin, err := h.chatService.IsUserAdminOfChat(context.Background(), chatID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to check membership",
		})
	}
	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "only admins can remove members",
		})
	}

	err = h.chatService.RemoveMember(context.Background(), chatID, memberID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to remove member",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "member removed successfully",
	})
}
