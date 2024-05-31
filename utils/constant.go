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
