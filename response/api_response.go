package response

import "time"

type ApiResponse[T any] struct {
	CorrelationID string    `json:"correlationid"`
	Error         *string   `json:"error"`
	Tin           time.Time `json:"tin"`
	Tout          time.Time `json:"tout"`
	Data          *T        `json:"data"`
	Success       bool      `json:"success"`
	StatusCode    int       `json:"statuscode"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statuscode"`
	Message    string `json:"message"`
	Error      any    `json:"error"`
}

func Success[T any](correlationID string, tin time.Time, statusCode int, data T) ApiResponse[T] {
	return SuccessPtr(correlationID, tin, statusCode, &data)
}

func SuccessPtr[T any](correlationID string, tin time.Time, statusCode int, data *T) ApiResponse[T] {
	return ApiResponse[T]{
		CorrelationID: correlationID,
		Error:         nil,
		Tin:           tin,
		Tout:          time.Now().UTC(),
		Data:          data,
		Success:       true,
		StatusCode:    statusCode,
	}
}

func Error[T any](correlationID string, tin time.Time, statusCode int, message string) ApiResponse[T] {
	return ApiResponse[T]{
		CorrelationID: correlationID,
		Error:         &message,
		Tin:           tin,
		Tout:          time.Now().UTC(),
		Data:          nil,
		Success:       false,
		StatusCode:    statusCode,
	}
}

func EmptySuccess(correlationID string, tin time.Time, statusCode int) ApiResponse[any] {
	return ApiResponse[any]{
		CorrelationID: correlationID,
		Error:         nil,
		Tin:           tin,
		Tout:          time.Now().UTC(),
		Data:          nil,
		Success:       true,
		StatusCode:    statusCode,
	}
}
