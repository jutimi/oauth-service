package utils

type key string

const (
	DATE_FORMAT       = "2006-01-02"
	DATE_TIME_FORMAT  = "2006-01-02 15:04:05"
	TIME_STAMP_FORMAT = "20060102150405"
)

const (
	GIN_CONTEXT_KEY       key = "GIN"
	USER_CONTEXT_KEY      key = "USER"
	SCOPE_CONTEXT_KEY     key = "SCOPE"
	WORKSPACE_CONTEXT_KEY key = "WORKSPACE"
)

const (
	ACCESS_TOKEN  = "access_token"
	REFRESH_TOKEN = "refresh_token"

	USER_AUTHORIZATION     = "User-Authorization"
	USER_ACCESS_TOKEN_IAT  = 15 * 60           // 15 minutes
	USER_REFRESH_TOKEN_IAT = 30 * 24 * 60 * 60 // 30 days

	WORKSPACE_AUTHORIZATION     = "Workspace-Authorization"
	Workspace_ACCESS_TOKEN_IAT  = 15 * 60           // 15 minutes
	Workspace_REFRESH_TOKEN_IAT = 30 * 24 * 60 * 60 // 30 days
)

const (
	DEBUG_MODE   = "debug"
	RELEASE_MODE = "release"
)

const (
	URL_TYPE_API = "api"
	URL_TYPE_CMS = "cms"
)

const (
	USER_SCOPE           = "user"
	WORKSPACE_SCOPE      = "workspace"
	USER_WORKSPACE_SCOPE = "user_workspace"
)
