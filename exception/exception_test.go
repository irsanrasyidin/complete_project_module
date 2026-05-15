package exception

import (
	"net/http"
	"testing"
)

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		code Code
		want int
	}{
		{InvalidArgumentCode, http.StatusBadRequest},
		{NotFoundCode, http.StatusNotFound},
		{AlreadyExistsCode, http.StatusConflict},
		{PermissionDeniedCode, http.StatusForbidden},
		{UnauthenticatedCode, http.StatusUnauthorized},
		{InternalErrorCode, http.StatusInternalServerError},
		{UnprocessableEntityCode, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		if got := HTTPStatus(tt.code); got != tt.want {
			t.Fatalf("HTTPStatus(%s) = %d, want %d", tt.code, got, tt.want)
		}
	}
}

func TestExceptionGetErrorPrefersDetails(t *testing.T) {
	exc := New(
		InvalidArgumentCode,
		"invalid request",
		WithError(assertErr("raw error")),
		WithDetails("public details"),
	)

	if got := exc.GetError(); got != "public details" {
		t.Fatalf("GetError() = %#v, want public details", got)
	}
}

type assertErr string

func (e assertErr) Error() string {
	return string(e)
}
