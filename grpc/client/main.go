package client_grpc

type ClientGRPCCollection struct {
	WSClient WorkspaceClient
}

func RegisterClientGRPC() ClientGRPCCollection {
	return ClientGRPCCollection{
		WSClient: NewWsClient(),
	}
}
