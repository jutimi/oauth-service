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

func GetScopeContext(ctx context.Context) (string, error) {
	ctxData := ctx.Value(SCOPE_CONTEXT_KEY)
	data, ok := ctxData.(string)
	if !ok {
		return "", errors.ErrUnsupported
	}

	return data, nil
}
