package rest

// MessageResponse модель успешного ответа
type MessageResponse struct {
	Message string `json:"message" example:"success"`
}

// ErrorResponse модель ответа об ошибке
type ErrorResponse struct {
	Message string `json:"message" example:"object not found"`
}
