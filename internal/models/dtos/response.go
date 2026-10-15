package dtos

import "time"

type ErrorResponse struct {
	Error error `json:"error"`
}

type Response[T any] struct {
	Status    bool      `json:"status"`
	Code      int       `json:"code,omitempty"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Data      T         `json:"data"`
}
