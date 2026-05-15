package exception

import "net/http"

type Code string

const (
	InvalidArgumentCode     Code = "INVALID_ARGUMENT"
	NotFoundCode            Code = "NOT_FOUND"
	AlreadyExistsCode       Code = "ALREADY_EXISTS"
	PermissionDeniedCode    Code = "PERMISSION_DENIED"
	UnauthenticatedCode     Code = "UNAUTHENTICATED"
	InternalErrorCode       Code = "INTERNAL"
	UnprocessableEntityCode Code = "UNPROCESSABLE_ENTITY"
)

func HTTPStatus(code Code) int {
	switch code {
	case InvalidArgumentCode:
		return http.StatusBadRequest
	case NotFoundCode:
		return http.StatusNotFound
	case AlreadyExistsCode:
		return http.StatusConflict
	case PermissionDeniedCode:
		return http.StatusForbidden
	case UnauthenticatedCode:
		return http.StatusUnauthorized
	case UnprocessableEntityCode:
		return http.StatusUnprocessableEntity
	case InternalErrorCode:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
