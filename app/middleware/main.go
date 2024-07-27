package middleware

type MiddlewareCollections struct {
	UserMW      Middleware
	WorkspaceMW Middleware
	OwnerMW     Middleware
}

func RegisterMiddleware() MiddlewareCollections {
	return MiddlewareCollections{
		UserMW:      NewUserMiddleware(),
		WorkspaceMW: NewWorkspaceMiddleware(),
		OwnerMW:     NewOwnerMiddleware(),
	}
}
