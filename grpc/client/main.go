package client_grpc

type ClientGRPCCollection struct {
	wsClient WorkspaceClient
}

func RegisterClientGRPC() ClientGRPCCollection {

	return ClientGRPCCollection{
		wsClient: NewWsClient(),
	}
}
