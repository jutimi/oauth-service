package utils

type key string

const (
	DATE_FORMAT       = "2006-01-02"
	DATE_TIME_FORMAT  = "2006-01-02 15:04:05"
	TIME_STAMP_FORMAT = "20060102150405"
)

const (
	GIN_CONTEXT_KEY  key = "GIN"
	USER_CONTEXT_KEY key = "USER"
)

const (
	USER_ACCESS_TOKEN_IAT  = 15 * 60           // 15 minutes
	USER_REFRESH_TOKEN_IAT = 30 * 24 * 60 * 60 // 30 days
)

const (
	DEBUG_MODE   = "debug"
	RELEASE_MODE = "release"
)
