package client

import (
	context "context"
	"log"
	"oauth-server/config"
	"oauth-server/package/errors"

	"github.com/jutimi/grpc-service/utils"
	"github.com/jutimi/grpc-service/workspace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type workspaceClient struct {
	conn                *grpc.ClientConn
	WorkspaceClient     workspace.WorkspaceRouteClient
	userWorkspaceClient workspace.UserWorkspaceRouteClient
}

type WorkspaceClient interface {
	GetWorkspaceByFilter(ctx context.Context, data *workspace.GetWorkspaceByFilterParams) (*workspace.WorkspaceResponse, error)
	GetUserWorkspaceByFilter(ctx context.Context, data *workspace.GetUserWorkspaceByFilterParams) (*workspace.UserWorkspaceResponse, error)
	CloseConn()
}

func NewWorkspaceClient() WorkspaceClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conf := config.GetConfiguration().GRPC

	// Connect to Workspace grpc server
	conn, err := grpc.NewClient(conf.WorkspaceUrl, opts...)
	if err != nil {
		log.Fatalf("Error connect to Workspace grpc server: %s", err.Error())
	}
	WorkspaceClient := workspace.NewWorkspaceRouteClient(conn)
	userWorkspaceClient := workspace.NewUserWorkspaceRouteClient(conn)

	return &workspaceClient{
		conn:                conn,
		WorkspaceClient:     WorkspaceClient,
		userWorkspaceClient: userWorkspaceClient,
	}
}

func (c *workspaceClient) GetWorkspaceByFilter(
	ctx context.Context,
	data *workspace.GetWorkspaceByFilterParams,
) (*workspace.WorkspaceResponse, error) {
	conf := config.GetConfiguration().Server
	ctx = utils.GenerateContext(utils.Metadata{
		Ctx:         ctx,
		ServiceName: conf.ServiceName,
	})

	resp, err := c.WorkspaceClient.GetWorkspaceByFilter(ctx, data)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if resp.Error != nil {
		return nil, errors.NewCustomError(int(resp.Error.ErrorCode), resp.Error.ErrorMessage)
	}
	if !resp.Success {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return resp, nil
}

func (c *workspaceClient) GetUserWorkspaceByFilter(
	ctx context.Context,
	data *workspace.GetUserWorkspaceByFilterParams,
) (*workspace.UserWorkspaceResponse, error) {
	conf := config.GetConfiguration().Server
	ctx = utils.GenerateContext(utils.Metadata{
		Ctx:         ctx,
		ServiceName: conf.ServiceName,
	})

	resp, err := c.userWorkspaceClient.GetUserWorkspaceByFilter(ctx, data)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if resp.Error != nil {
		return nil, errors.NewCustomError(int(resp.Error.ErrorCode), resp.Error.ErrorMessage)
	}
	if !resp.Success {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return resp, nil
}

func (c *workspaceClient) CloseConn() {
	c.conn.Close()
}
