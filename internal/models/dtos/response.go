package dtos

type Response[T any] struct {
	Code    string `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Error   *error `json:"error,omitempty"`
}
