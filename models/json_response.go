package models

type SuccessPayload[T any] struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    T      `json:"data"`
}

type FailPayload struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
