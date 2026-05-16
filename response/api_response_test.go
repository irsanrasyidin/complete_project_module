package response

import (
	"net/http"
	"testing"
	"time"

	"github.com/irsanrasyidin/complete_project_module/exception"
)

func TestSuccess(t *testing.T) {
	tin := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)

	res := Success("corr-1", tin, http.StatusOK, "ok")

	if !res.Success {
		t.Fatal("expected success response")
	}
	if res.Error != nil {
		t.Fatal("expected nil error")
	}
	if res.Data == nil || *res.Data != "ok" {
		t.Fatalf("expected data ok, got %#v", res.Data)
	}
	if res.CorrelationID != "corr-1" {
		t.Fatalf("unexpected correlation id: %s", res.CorrelationID)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", res.StatusCode)
	}
}

func TestError(t *testing.T) {
	tin := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC)

	res := Error[any]("corr-1", tin, http.StatusBadRequest, "invalid payload")

	if res.Success {
		t.Fatal("expected error response")
	}
	if res.Error == nil || *res.Error != "invalid payload" {
		t.Fatalf("unexpected error: %#v", res.Error)
	}
	if res.Data != nil {
		t.Fatalf("expected nil data, got %#v", res.Data)
	}
}

func TestNewExceptionErrorResponseHidesInternalError(t *testing.T) {
	exc := exception.New(
		exception.InternalErrorCode,
		"internal server error",
		exception.WithError(assertErr("database password leaked")),
	)

	res := NewExceptionErrorResponse(exc)

	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("unexpected status code: %d", res.StatusCode)
	}
	if res.Error != nil {
		t.Fatalf("expected hidden internal error, got %#v", res.Error)
	}
}

func TestNewExceptionErrorResponseShowsNonInternalError(t *testing.T) {
	exc := exception.New(
		exception.InvalidArgumentCode,
		"invalid request",
		exception.WithDetails(map[string]string{"amount": "required"}),
	)

	res := NewExceptionErrorResponse(exc)

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", res.StatusCode)
	}
	if res.Error == nil {
		t.Fatal("expected public error details")
	}
}

type assertErr string

func (e assertErr) Error() string {
	return string(e)
}
