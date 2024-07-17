package gRPC

import "github.com/jutimi/grpc-service/oauth"

type GRPCServer interface {
	oauth.UserRouteServer
	oauth.OAuthRouteServer
}

type GRPCClient interface {
}
