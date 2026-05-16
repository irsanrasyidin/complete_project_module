package sso

import (
	"context"
	"net/http"
)

type contextKey string

const (
	CorrelationIDKey contextKey = "sso-correlation-id"
	TenantCodeKey    contextKey = "sso-tenant-code"
	RealIPKey        contextKey = "sso-real-ip"
	ForwardedForKey  contextKey = "sso-forwarded-for"
	UserNameKey      contextKey = "sso-user-name"
	OriginURLKey     contextKey = "sso-origin-url"
	RolesKey         contextKey = "sso-roles"
	TokenIDKey       contextKey = "sso-token-id"
	SalesCodeKey     contextKey = "sso-sales-code"
	BranchKey        contextKey = "sso-branch"

	contextProcessKey contextKey = "sso-context-process"
	requestIDKey      contextKey = "sso-request-id"
	requestIPKey      contextKey = "sso-request-ip"
	hostKey           contextKey = "sso-host"
	baseURLKey        contextKey = "sso-base-url"
	langKey           contextKey = "sso-lang"
	userIDKey         contextKey = "sso-user-id"
	usernameKey       contextKey = "sso-username"
	emailKey          contextKey = "sso-email"
	identityRolesKey  contextKey = "sso-identity-roles"
)

const (
	CorrelationIDHeader = "X-Correlation-ID"
	TenantCodeHeader    = "X-Sso-Tenantcode"
	RealIPHeader        = "X-Real-Ip"
	ForwardedForHeader  = "X-Forwarded-For"
	UserNameHeader      = "X-Sso-Username"
	OriginURLHeader     = "X-Origin-Url"
	RolesHeader         = "X-Sso-Roles"
	TokenIDHeader       = "X-Sso-Tokenid"
	SalesCodeHeader     = "X-Sso-Salescode"
	BranchHeader        = "X-Sso-Branch"
)

type AuditInfo struct {
	CorrelationID string
	TenantCode    string
	RealIP        string
	ForwardedFor  string
	UserName      string
	OriginURL     string
	Roles         string
	TokenID       string
	SalesCode     string
	Branch        string
}

type Identity struct {
	UserID   string
	Username string
	Email    string
	Roles    []string
}

func AuditInfoFromHeader(header http.Header) AuditInfo {
	return AuditInfo{
		CorrelationID: header.Get(CorrelationIDHeader),
		TenantCode:    header.Get(TenantCodeHeader),
		RealIP:        header.Get(RealIPHeader),
		ForwardedFor:  header.Get(ForwardedForHeader),
		UserName:      header.Get(UserNameHeader),
		OriginURL:     header.Get(OriginURLHeader),
		Roles:         header.Get(RolesHeader),
		TokenID:       header.Get(TokenIDHeader),
		SalesCode:     header.Get(SalesCodeHeader),
		Branch:        header.Get(BranchHeader),
	}
}

func SetAuditInfo(ctx context.Context, auditInfo AuditInfo) context.Context {
	ctx = SetCorrelationID(ctx, auditInfo.CorrelationID)
	ctx = SetTenantCode(ctx, auditInfo.TenantCode)
	ctx = SetRealIP(ctx, auditInfo.RealIP)
	ctx = SetForwardedFor(ctx, auditInfo.ForwardedFor)
	ctx = SetUserName(ctx, auditInfo.UserName)
	ctx = SetOriginURL(ctx, auditInfo.OriginURL)
	ctx = SetRolesRaw(ctx, auditInfo.Roles)
	ctx = SetTokenID(ctx, auditInfo.TokenID)
	ctx = SetSalesCode(ctx, auditInfo.SalesCode)
	return SetBranch(ctx, auditInfo.Branch)
}

func AuditInfoFromContext(ctx context.Context) AuditInfo {
	return AuditInfo{
		CorrelationID: GetCorrelationID(ctx),
		TenantCode:    GetTenantCode(ctx),
		RealIP:        GetRealIP(ctx),
		ForwardedFor:  GetForwardedFor(ctx),
		UserName:      GetUserName(ctx),
		OriginURL:     GetOriginURL(ctx),
		Roles:         GetRolesRaw(ctx),
		TokenID:       GetTokenID(ctx),
		SalesCode:     GetSalesCode(ctx),
		Branch:        GetBranch(ctx),
	}
}

func ContextWithAuditInfoFromHeader(ctx context.Context, header http.Header) context.Context {
	return SetAuditInfo(ctx, AuditInfoFromHeader(header))
}

func SetCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, correlationID)
}

func GetCorrelationID(ctx context.Context) string {
	return stringValue(ctx, CorrelationIDKey)
}

func SetTenantCode(ctx context.Context, tenantCode string) context.Context {
	return context.WithValue(ctx, TenantCodeKey, tenantCode)
}

func GetTenantCode(ctx context.Context) string {
	return stringValue(ctx, TenantCodeKey)
}

func SetRealIP(ctx context.Context, realIP string) context.Context {
	return context.WithValue(ctx, RealIPKey, realIP)
}

func GetRealIP(ctx context.Context) string {
	return stringValue(ctx, RealIPKey)
}

func SetForwardedFor(ctx context.Context, forwardedFor string) context.Context {
	return context.WithValue(ctx, ForwardedForKey, forwardedFor)
}

func GetForwardedFor(ctx context.Context) string {
	return stringValue(ctx, ForwardedForKey)
}

func SetUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, UserNameKey, userName)
}

func GetUserName(ctx context.Context) string {
	return stringValue(ctx, UserNameKey)
}

func SetOriginURL(ctx context.Context, originURL string) context.Context {
	return context.WithValue(ctx, OriginURLKey, originURL)
}

func GetOriginURL(ctx context.Context) string {
	return stringValue(ctx, OriginURLKey)
}

func SetRolesRaw(ctx context.Context, roles string) context.Context {
	return context.WithValue(ctx, RolesKey, roles)
}

func GetRolesRaw(ctx context.Context) string {
	return stringValue(ctx, RolesKey)
}

func SetTokenID(ctx context.Context, tokenID string) context.Context {
	return context.WithValue(ctx, TokenIDKey, tokenID)
}

func GetTokenID(ctx context.Context) string {
	return stringValue(ctx, TokenIDKey)
}

func SetSalesCode(ctx context.Context, salesCode string) context.Context {
	return context.WithValue(ctx, SalesCodeKey, salesCode)
}

func GetSalesCode(ctx context.Context) string {
	return stringValue(ctx, SalesCodeKey)
}

func SetBranch(ctx context.Context, branch string) context.Context {
	return context.WithValue(ctx, BranchKey, branch)
}

func GetBranch(ctx context.Context) string {
	return stringValue(ctx, BranchKey)
}

func SetContextProcess(ctx context.Context, process string) context.Context {
	return context.WithValue(ctx, contextProcessKey, process)
}

func GetContextProcess(ctx context.Context) string {
	return stringValue(ctx, contextProcessKey)
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	return stringValue(ctx, requestIDKey)
}

func SetRequestIP(ctx context.Context, requestIP string) context.Context {
	return context.WithValue(ctx, requestIPKey, requestIP)
}

func GetRequestIP(ctx context.Context) string {
	return stringValue(ctx, requestIPKey)
}

func SetHost(ctx context.Context, host string) context.Context {
	return context.WithValue(ctx, hostKey, host)
}

func GetHost(ctx context.Context) string {
	return stringValue(ctx, hostKey)
}

func SetBaseURL(ctx context.Context, baseURL string) context.Context {
	return context.WithValue(ctx, baseURLKey, baseURL)
}

func GetBaseURL(ctx context.Context) string {
	return stringValue(ctx, baseURLKey)
}

func SetLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langKey, lang)
}

func GetLang(ctx context.Context) string {
	return stringValue(ctx, langKey)
}

func SetIdentity(ctx context.Context, identity Identity) context.Context {
	ctx = context.WithValue(ctx, userIDKey, identity.UserID)
	ctx = context.WithValue(ctx, usernameKey, identity.Username)
	ctx = context.WithValue(ctx, emailKey, identity.Email)
	return context.WithValue(ctx, identityRolesKey, identity.Roles)
}

func GetIdentity(ctx context.Context) Identity {
	return Identity{
		UserID:   GetUserID(ctx),
		Username: GetUsername(ctx),
		Email:    GetEmail(ctx),
		Roles:    GetRoles(ctx),
	}
}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) string {
	return stringValue(ctx, userIDKey)
}

func SetUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func GetUsername(ctx context.Context) string {
	return stringValue(ctx, usernameKey)
}

func SetEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailKey, email)
}

func GetEmail(ctx context.Context) string {
	return stringValue(ctx, emailKey)
}

func SetRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, identityRolesKey, roles)
}

func GetRoles(ctx context.Context) []string {
	roles, ok := ctx.Value(identityRolesKey).([]string)
	if !ok {
		return nil
	}
	return roles
}

func stringValue(ctx context.Context, key contextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}
