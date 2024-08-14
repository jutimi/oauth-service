package utils

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
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

func GetScopeContext[T string](ctx context.Context, key key) (T, error) {
	ctxData := ctx.Value(key)
	data, ok := ctxData.(T)
	if !ok {
		return data, errors.ErrUnsupported
	}

	return data, nil
}

func GetGinContext[T *UserPayload | *WorkspacePayload](ctx context.Context, key string) (T, error) {
	ginCtx := ctx.Value(GIN_CONTEXT_KEY).(*gin.Context)
	ctxData, ok := ginCtx.Get(key)
	if !ok {
		return nil, errors.ErrUnsupported
	}

	data, ok := ctxData.(T)
	if !ok {
		return nil, fmt.Errorf("invalid payload")
	}

	return data, nil
}
