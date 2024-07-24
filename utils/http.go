package utils

import (
	"context"
	"errors"
)

func FormatSuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"success": true,
		"data":    data,
	}
}

func FormatErrorResponse(err error) map[string]interface{} {
	return map[string]interface{}{
		"success": false,
		"error":   err,
	}
}

func GetScopeContext[T string | *UserPayload | *WorkspacePayload](ctx context.Context, key string) (T, error) {
	ctxData := ctx.Value(key)
	data, ok := ctxData.(T)
	if !ok {
		return data, errors.ErrUnsupported
	}

	return data, nil
}
