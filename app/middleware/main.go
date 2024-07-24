package middleware

type MiddlewareCollections struct {
	UserMW      Middleware
	WorkspaceMW Middleware
}

func RegisterMiddleware() MiddlewareCollections {
	return MiddlewareCollections{
		UserMW:      NewUserMiddleware(),
		WorkspaceMW: NewWorkspaceMiddleware(),
	}
}
