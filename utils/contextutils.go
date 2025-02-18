package utils

import (
	"context"
)

const (
	ContextProcess = "context-process"
	RequestID      = "request-id"
	RequestIP      = "request-ip"
	Host           = "host"
	BaseURL        = "base-url"
	Lang           = "lang"
)

func SetContextProcess(ctx context.Context, process string) context.Context {
	return context.WithValue(ctx, ContextProcess, process)
}

func GetContextProcess(ctx context.Context) string {
	process, ok := ctx.Value(ContextProcess).(string)
	if !ok {
		return ""
	}
	return process
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestID).(string)
	if !ok {
		return ""
	}
	return requestID
}

func GetRequestIP(ctx context.Context) string {
	requestIP, ok := ctx.Value(RequestIP).(string)
	if !ok {
		return ""
	}
	return requestIP
}

func GetHost(ctx context.Context) string {
	host, ok := ctx.Value(Host).(string)
	if !ok {
		return ""
	}
	return host
}

func GetBaseURL(ctx context.Context) string {
	baseURL, ok := ctx.Value(BaseURL).(string)
	if !ok {
		return ""
	}
	return baseURL
}

func GetLang(ctx context.Context) string {
	lang, ok := ctx.Value(Lang).(string)
	if !ok {
		return ""
	}
	return lang
}
