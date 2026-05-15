package response

import "github.com/irsanrasyidin/complete_project/module/exception"

type JSONAborter interface {
	AbortWithStatusJSON(code int, jsonObj any)
}

func NewExceptionErrorResponse(exc *exception.Exception) *ErrorResponse {
	if exc == nil {
		exc = exception.New(exception.InternalErrorCode, "internal server error")
	}

	var err any
	if exc.Code != exception.InternalErrorCode {
		err = exc.GetError()
	}

	return &ErrorResponse{
		StatusCode: exc.GetHttpCode(),
		Message:    exc.Message,
		Error:      err,
	}
}

func ExceptionJSON(c JSONAborter, exc *exception.Exception) {
	errRes := NewExceptionErrorResponse(exc)
	c.AbortWithStatusJSON(errRes.StatusCode, errRes)
}
