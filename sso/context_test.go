package sso

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestRequestContextValues(t *testing.T) {
	ctx := context.Background()
	ctx = SetRequestID(ctx, "req-1")
	ctx = SetRequestIP(ctx, "127.0.0.1")
	ctx = SetHost(ctx, "example.test")
	ctx = SetBaseURL(ctx, "https://example.test")
	ctx = SetLang(ctx, "id")
	ctx = SetContextProcess(ctx, "login")

	if GetRequestID(ctx) != "req-1" {
		t.Fatalf("unexpected request id: %s", GetRequestID(ctx))
	}
	if GetRequestIP(ctx) != "127.0.0.1" {
		t.Fatalf("unexpected request ip: %s", GetRequestIP(ctx))
	}
	if GetHost(ctx) != "example.test" {
		t.Fatalf("unexpected host: %s", GetHost(ctx))
	}
	if GetBaseURL(ctx) != "https://example.test" {
		t.Fatalf("unexpected base url: %s", GetBaseURL(ctx))
	}
	if GetLang(ctx) != "id" {
		t.Fatalf("unexpected lang: %s", GetLang(ctx))
	}
	if GetContextProcess(ctx) != "login" {
		t.Fatalf("unexpected context process: %s", GetContextProcess(ctx))
	}
}

func TestIdentityContextValues(t *testing.T) {
	ctx := SetIdentity(context.Background(), Identity{
		UserID:   "usr-1",
		Username: "irsan",
		Email:    "irsan@example.test",
		Roles:    []string{"admin", "user"},
	})

	identity := GetIdentity(ctx)

	if identity.UserID != "usr-1" {
		t.Fatalf("unexpected user id: %s", identity.UserID)
	}
	if identity.Username != "irsan" {
		t.Fatalf("unexpected username: %s", identity.Username)
	}
	if identity.Email != "irsan@example.test" {
		t.Fatalf("unexpected email: %s", identity.Email)
	}
	if !reflect.DeepEqual(identity.Roles, []string{"admin", "user"}) {
		t.Fatalf("unexpected roles: %#v", identity.Roles)
	}
}

func TestAuditInfoFromHeader(t *testing.T) {
	header := http.Header{}
	header.Set(CorrelationIDHeader, "corr-1")
	header.Set(TenantCodeHeader, "cimb")
	header.Set(RealIPHeader, "127.0.0.1")
	header.Set(ForwardedForHeader, "*")
	header.Set(UserNameHeader, "seiya")
	header.Set(OriginURLHeader, "https://example.test")
	header.Set(RolesHeader, "admin")
	header.Set(TokenIDHeader, "token-1")
	header.Set(SalesCodeHeader, "sales-1")
	header.Set(BranchHeader, "branch-1")

	auditInfo := AuditInfoFromHeader(header)

	if auditInfo.CorrelationID != "corr-1" {
		t.Fatalf("unexpected correlation id: %s", auditInfo.CorrelationID)
	}
	if auditInfo.TenantCode != "cimb" {
		t.Fatalf("unexpected tenant code: %s", auditInfo.TenantCode)
	}
	if auditInfo.RealIP != "127.0.0.1" {
		t.Fatalf("unexpected real ip: %s", auditInfo.RealIP)
	}
	if auditInfo.ForwardedFor != "*" {
		t.Fatalf("unexpected forwarded for: %s", auditInfo.ForwardedFor)
	}
	if auditInfo.UserName != "seiya" {
		t.Fatalf("unexpected username: %s", auditInfo.UserName)
	}
	if auditInfo.Roles != "admin" {
		t.Fatalf("unexpected roles: %s", auditInfo.Roles)
	}
}

func TestAuditInfoContextValues(t *testing.T) {
	ctx := SetAuditInfo(context.Background(), AuditInfo{
		CorrelationID: "corr-1",
		TenantCode:    "cimb",
		RealIP:        "127.0.0.1",
		ForwardedFor:  "*",
		UserName:      "seiya",
		OriginURL:     "https://example.test",
		Roles:         "admin",
		TokenID:       "token-1",
		SalesCode:     "sales-1",
		Branch:        "branch-1",
	})

	auditInfo := AuditInfoFromContext(ctx)

	if auditInfo.CorrelationID != "corr-1" {
		t.Fatalf("unexpected correlation id: %s", auditInfo.CorrelationID)
	}
	if auditInfo.TenantCode != "cimb" {
		t.Fatalf("unexpected tenant code: %s", auditInfo.TenantCode)
	}
	if auditInfo.RealIP != "127.0.0.1" {
		t.Fatalf("unexpected real ip: %s", auditInfo.RealIP)
	}
	if auditInfo.UserName != "seiya" {
		t.Fatalf("unexpected username: %s", auditInfo.UserName)
	}
	if auditInfo.Roles != "admin" {
		t.Fatalf("unexpected roles: %s", auditInfo.Roles)
	}
}

func TestEmptyContextValues(t *testing.T) {
	ctx := context.Background()

	if GetRequestID(ctx) != "" {
		t.Fatalf("expected empty request id")
	}
	identity := GetIdentity(ctx)
	if identity.UserID != "" || identity.Username != "" || identity.Email != "" || identity.Roles != nil {
		t.Fatalf("expected empty identity, got %#v", identity)
	}
}
