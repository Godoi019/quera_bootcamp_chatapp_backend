package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

var (
	ErrRequiredBodyNotFound       = errors.New("missing required body data")
	ErrRequiredUrlQueriesNotFound = errors.New("missing required url queries")
)

// customErrorHandler handles errors in a unified way
func CustomErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}

func RespondError(c fiber.Ctx, status int, err error, messages interface{}) error {
	resp := APIResponse{
		Success:  false,
		Messages: messages,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	return c.Status(status).JSON(resp)
}
