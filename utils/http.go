package utils

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
