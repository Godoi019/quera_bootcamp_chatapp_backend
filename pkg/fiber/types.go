package fiber

type ServiceGetter[T any] func() T

type APIResponse struct {
	Success  bool        `json:"success"`
	Messages interface{} `json:"messages,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Error    string      `json:"error,omitempty"`
}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

type ValidationError struct {
	Errors  []*ErrorResponse
	Message error
}

type ctxKey string

const fiberCtxKey ctxKey = "fiberCtx"
