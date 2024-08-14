package server

import (
	context "context"
	"oauth-server/app/helper"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/config"
	"oauth-server/external/client"
	"oauth-server/package/errors"
	"oauth-server/utils"
	"strings"

	"github.com/jutimi/grpc-service/oauth"
	grpc_utils "github.com/jutimi/grpc-service/utils"
	"github.com/jutimi/grpc-service/workspace"

	"github.com/google/uuid"
)

type grpcServer struct {
	oauth.UnimplementedOAuthRouteServer
	oauth.UnimplementedUserRouteServer

	postgresRepo postgres_repository.PostgresRepositoryCollections
	helper       helper.HelperCollections
}

type OAuthServer interface {
	oauth.UserRouteServer
	oauth.OAuthRouteServer
}

func NewGRPCServer(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	helper helper.HelperCollections,
) OAuthServer {
	return &grpcServer{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}

func (s *grpcServer) GetUsersByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UsersResponse, error) {

	var usersRes []*oauth.UserDetail

	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		customErr := errors.New(errors.ErrCodeInternalServerError)
		return &oauth.UsersResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	users, err := s.postgresRepo.UserRepo.FindByFilter(ctx, filter)
	if err != nil {
		customErr := errors.New(errors.ErrCodeInternalServerError)
		return &oauth.UsersResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	for _, user := range users {
		usersRes = append(usersRes, &oauth.UserDetail{
			Id:             user.Id.String(),
			PhoneNumber:    user.PhoneNumber,
			Email:          user.Email,
			IsActive:       user.IsActive,
			LimitWorkspace: int32(user.LimitWorkspace),
		})
	}

	return &oauth.UsersResponse{
		Success: true,
		Data:    usersRes,
		Error:   nil,
	}, nil
}

func (s *grpcServer) GetUserByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UserResponse, error) {
	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		customErr := errors.New(errors.ErrCodeInternalServerError)
		return &oauth.UserResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, filter)
	if err != nil {
		customErr := errors.New(errors.ErrCodeUserNotFound)
		return &oauth.UserResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	return &oauth.UserResponse{
		Success: true,
		Data: &oauth.UserDetail{
			Id:             user.Id.String(),
			PhoneNumber:    user.PhoneNumber,
			Email:          user.Email,
			IsActive:       user.IsActive,
			LimitWorkspace: int32(user.LimitWorkspace),
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) CreateUser(ctx context.Context, data *oauth.CreateUserParams) (*oauth.UserResponse, error) {
	user, err := s.helper.UserHelper.CreateUser(ctx, &helper.CreateUserParams{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
		Password:    data.Password,
	})
	if err != nil {
		cusErr := err.(*errors.CustomError)
		return &oauth.UserResponse{
			Success: false,
			Error:   grpc_utils.FormatErrorResponse(int32(cusErr.GetCode()), cusErr.Error()),
			Data:    nil,
		}, err
	}

	return &oauth.UserResponse{
		Success: true,
		Error:   nil,
		Data: &oauth.UserDetail{
			Id:             user.Id.String(),
			PhoneNumber:    user.PhoneNumber,
			Email:          user.Email,
			IsActive:       user.IsActive,
			LimitWorkspace: int32(user.LimitWorkspace),
		},
	}, nil
}

func (s *grpcServer) VerifyUserToken(ctx context.Context, data *oauth.VerifyTokenParams) (*oauth.VerifyTokenResponse, error) {
	conf := config.GetConfiguration().Jwt
	customErr := errors.New(errors.ErrCodeInternalServerError)

	payload, err := utils.VerifyUserToken(data.Token, conf.UserAccessTokenKey)
	if err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	isActive := true
	if _, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, &repository.FindUserByFilter{
		Id:       &payload.Id,
		IsActive: &isActive,
	}); err != nil {
		customErr = errors.New(errors.ErrCodeUserNotFound)
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	return &oauth.VerifyTokenResponse{Success: true}, nil
}

func (s *grpcServer) VerifyWorkspaceToken(ctx context.Context, data *oauth.VerifyTokenParams) (*oauth.VerifyTokenResponse, error) {
	conf := config.GetConfiguration().Jwt
	customErr := errors.New(errors.ErrCodeInternalServerError)

	payload, err := utils.VerifyWorkspaceToken(data.Token, conf.WorkspaceAccessTokenKey)
	if err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	userWorkspaceId := payload.UserWorkspaceId.String()
	WorkspaceId := payload.WorkspaceId.String()
	isActive := true
	WorkspaceGRPC := client.NewWorkspaceClient()
	defer WorkspaceGRPC.CloseConn()

	// Check workspace data
	Workspace, err := WorkspaceGRPC.GetWorkspaceByFilter(ctx, &workspace.GetWorkspaceByFilterParams{
		Id:       &WorkspaceId,
		IsActive: &isActive,
	})
	if err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}
	if Workspace.Error != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error:   Workspace.Error,
		}, nil
	}
	// Check user workspace data
	userWorkspace, err := WorkspaceGRPC.GetUserWorkspaceByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
		Id:       &userWorkspaceId,
		IsActive: &isActive,
	})
	if err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}
	if userWorkspace.Error != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error:   userWorkspace.Error,
		}, nil
	}

	return &oauth.VerifyTokenResponse{Success: true}, nil
}

func (s *grpcServer) VerifyWorkspacePermission(ctx context.Context, data *oauth.VerifyPermissionParams) (*oauth.VerifyTokenResponse, error) {
	conf := config.GetConfiguration().Jwt
	customErr := errors.New(errors.ErrCodeForbidden)

	payload, err := utils.VerifyWorkspaceToken(data.Token, conf.WorkspaceAccessTokenKey)
	if err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	url := strings.Split(data.GetUrl(), "/")
	urlType := url[3]
	resource := url[5]
	action := url[6]
	if urlType == utils.URL_TYPE_API {
		return &oauth.VerifyTokenResponse{Success: true}, nil
	}

	permission, err := s.helper.PermissionHelper.GetURLPermission(ctx, resource, action)
	if err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(err.(*errors.CustomError).GetCode()),
				err.(*errors.CustomError).Error(),
			),
		}, nil
	}

	// Check user permission
	if _, err := s.postgresRepo.PermissionRepo.FindOneByFilter(ctx, &repository.FindPermissionByFilter{
		UserWorkspaceId: &payload.UserWorkspaceId,
		Permission:      &permission,
	}); err != nil {
		return &oauth.VerifyTokenResponse{
			Success: false,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	return &oauth.VerifyTokenResponse{Success: true}, nil
}

func (s *grpcServer) BulkCreateUsers(ctx context.Context, data *oauth.CreateUsersParams) (*oauth.UsersResponse, error) {
	params := make([]helper.CreateUserParams, 0)
	usersRes := make([]*oauth.UserDetail, 0)

	for _, user := range data.Users {
		params = append(params, helper.CreateUserParams{
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Password:    user.Password,
		})
	}

	users, err := s.helper.UserHelper.CreateUsers(ctx, params)
	if err != nil {
		cusErr := err.(*errors.CustomError)
		return &oauth.UsersResponse{
			Success: false,
			Error:   grpc_utils.FormatErrorResponse(int32(cusErr.GetCode()), cusErr.Error()),
			Data:    nil,
		}, err
	}
	for _, user := range users {
		usersRes = append(usersRes, &oauth.UserDetail{
			Id:             user.Id.String(),
			PhoneNumber:    user.PhoneNumber,
			Email:          user.Email,
			IsActive:       user.IsActive,
			LimitWorkspace: int32(user.LimitWorkspace),
		})
	}

	return &oauth.UsersResponse{
		Success: true,
		Error:   nil,
		Data:    usersRes,
	}, nil
}

// ------------------------ Helper ------------------------
func convertUserParamsToFilter(data *oauth.GetUserByFilterParams) (*repository.FindUserByFilter, error) {
	var userId uuid.UUID
	var userIds []uuid.UUID
	var err error
	var limit, offset int

	if data.Limit != nil {
		limit = int(*data.Limit)
	}
	if data.Offset != nil {
		offset = int(*data.Offset)
	}

	if data.Id != nil {
		userId, err = utils.ConvertStringToUUID(*data.Id)
		if err != nil {
			return nil, err
		}
	}
	if data.Ids != nil {
		for _, id := range data.Ids {
			userId, err = utils.ConvertStringToUUID(id)
			if err != nil {
				return nil, err
			}

			userIds = append(userIds, userId)
		}
	}

	return &repository.FindUserByFilter{
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		Id:           &userId,
		Limit:        &limit,
		Offset:       &offset,
		Ids:          userIds,
		Emails:       data.Emails,
		PhoneNumbers: data.PhoneNumbers,
		IsActive:     data.IsActive,
	}, nil
}
