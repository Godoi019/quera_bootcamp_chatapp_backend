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

type MessageHandler struct {
	messageService *service.MessageService
	chatService    *service.ChatService
}

func NewMessageHandler(client *ent.Client) *MessageHandler {
	return &MessageHandler{
		messageService: service.NewMessageService(client),
		chatService:    service.NewChatService(client),
	}
}

func (h *MessageHandler) SendMessage(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	req := new(model.SendMessageRequest)
	if err := f.ParseRequestBody(c, req); err != nil {
		return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
	}

	// Check if user is a member of the chat
	isMember, err := h.chatService.IsUserMemberOfChat(context.Background(), req.ChatID, userID)
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

	// Send message
	newMessage, err := h.messageService.SendMessage(context.Background(), req.ChatID, userID, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to send message",
		})
	}

	senderID := userID
	chatID := req.ChatID
	if newMessage.Edges.Sender != nil {
		senderID = newMessage.Edges.Sender.ID
	}
	if newMessage.Edges.Chat != nil {
		chatID = newMessage.Edges.Chat.ID
	}

	return c.Status(fiber.StatusCreated).JSON(model.MessageResponse{
		ID:        newMessage.ID,
		Content:   newMessage.Content,
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: newMessage.CreatedAt,
		UpdatedAt: newMessage.UpdatedAt,
	})
}

func (h *MessageHandler) GetMessage(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	messageID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid message id",
		})
	}

	// Get message
	msg, err := h.messageService.GetMessageByID(context.Background(), messageID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.ErrorResponse{
			Error: "message not found",
		})
	}

	// Get IDs from edges
	senderID := 0
	chatID := 0
	if msg.Edges.Sender != nil {
		senderID = msg.Edges.Sender.ID
	}
	if msg.Edges.Chat != nil {
		chatID = msg.Edges.Chat.ID
	}

	// Check if user is a member of the chat
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

	response := model.MessageResponse{
		ID:        msg.ID,
		Content:   msg.Content,
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	}

	if msg.Edges.Sender != nil {
		response.Sender = &model.UserProfile{
			ID:          msg.Edges.Sender.ID,
			Username:    msg.Edges.Sender.Username,
			DisplayName: msg.Edges.Sender.DisplayName,
			CreatedAt:   msg.Edges.Sender.CreatedAt,
			LastSeen:    msg.Edges.Sender.LastSeen,
		}
	}

	return c.JSON(response)
}

func (h *MessageHandler) ListMessages(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	chatID, err := utils.ParamsInt(c, "chatId")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid chat id",
		})
	}

	limit := utils.QueryInt(c, "limit", 50)
	offset := utils.QueryInt(c, "offset", 0)

	// Check if user is a member of the chat
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

	messages, err := h.messageService.ListChatMessages(context.Background(), chatID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to list messages",
		})
	}

	messageResponses := make([]model.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		senderID := 0
		chatID := 0
		if msg.Edges.Sender != nil {
			senderID = msg.Edges.Sender.ID
		}
		if msg.Edges.Chat != nil {
			chatID = msg.Edges.Chat.ID
		}

		response := model.MessageResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			SenderID:  senderID,
			ChatID:    chatID,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		}

		if msg.Edges.Sender != nil {
			response.Sender = &model.UserProfile{
				ID:          msg.Edges.Sender.ID,
				Username:    msg.Edges.Sender.Username,
				DisplayName: msg.Edges.Sender.DisplayName,
				CreatedAt:   msg.Edges.Sender.CreatedAt,
				LastSeen:    msg.Edges.Sender.LastSeen,
			}
		}

		messageResponses = append(messageResponses, response)
	}

	return c.JSON(messageResponses)
}

func (h *MessageHandler) UpdateMessage(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	messageID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid message id",
		})
	}

	req := new(model.UpdateMessageRequest)
	if err := f.ParseRequestBody(c, req); err != nil {
		return f.RespondError(c, fiber.StatusBadRequest, err.Message, err.Errors)
	}

	// Check if user is the sender
	isSender, err := h.messageService.IsUserSenderOfMessage(context.Background(), messageID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to check message ownership",
		})
	}
	if !isSender {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "you can only edit your own messages",
		})
	}

	// Update message
	updatedMessage, err := h.messageService.UpdateMessage(context.Background(), messageID, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to update message",
		})
	}

	// Reload with edges
	updatedMessage, _ = h.messageService.GetMessageByID(context.Background(), messageID)

	senderID := userID
	chatID := 0
	if updatedMessage.Edges.Sender != nil {
		senderID = updatedMessage.Edges.Sender.ID
	}
	if updatedMessage.Edges.Chat != nil {
		chatID = updatedMessage.Edges.Chat.ID
	}

	return c.JSON(model.MessageResponse{
		ID:        updatedMessage.ID,
		Content:   updatedMessage.Content,
		SenderID:  senderID,
		ChatID:    chatID,
		CreatedAt: updatedMessage.CreatedAt,
		UpdatedAt: updatedMessage.UpdatedAt,
	})
}

func (h *MessageHandler) DeleteMessage(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	messageID, err := utils.ParamsInt(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error: "invalid message id",
		})
	}

	// Check if user is the sender
	isSender, err := h.messageService.IsUserSenderOfMessage(context.Background(), messageID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to check message ownership",
		})
	}
	if !isSender {
		return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
			Error: "you can only delete your own messages",
		})
	}

	// Delete message
	err = h.messageService.DeleteMessage(context.Background(), messageID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error: "failed to delete message",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
