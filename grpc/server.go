package server_grpc

import (
	context "context"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/package/errors"
	"oauth-server/utils"

	"github.com/jutimi/grpc-service/common"
	"github.com/jutimi/grpc-service/oauth"
	grpc_utils "github.com/jutimi/grpc-service/utils"

	"github.com/google/uuid"
)

type grpcServer struct {
	oauth.UnimplementedOAuthRouteServer
	oauth.UnimplementedUserRouteServer

	postgresRepo postgres_repository.PostgresRepositoryCollections
}

type OAuthServer interface {
	oauth.UserRouteServer
	oauth.OAuthRouteServer
}

func NewGRPCServer(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) OAuthServer {
	return &grpcServer{
		postgresRepo: postgresRepo,
	}
}

func (s *grpcServer) GetUserById(ctx context.Context, data *common.GetByIdParams) (*oauth.UserResponse, error) {
	userId, err := utils.ConvertStringToUUID(data.Id)
	if err != nil {
		customErr := errors.NewCustomError(errors.ErrCodeInternalServerError, err.Error())
		return &oauth.UserResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, &repository.FindUserByFilter{
		ID: &userId,
	})
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
			Id:          user.ID.String(),
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) GetUsersByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UsersResponse, error) {

	var usersRes []*oauth.UserDetail

	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		customErr := errors.NewCustomError(errors.ErrCodeInternalServerError, err.Error())
		return &oauth.UsersResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	users, err := s.postgresRepo.UserRepo.FindByFilter(ctx, nil, filter)
	if err != nil {
		customErr := errors.NewCustomError(errors.ErrCodeInternalServerError, err.Error())
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
			Id:          user.ID.String(),
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
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
		customErr := errors.NewCustomError(errors.ErrCodeInternalServerError, err.Error())
		return &oauth.UserResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	user, err := s.postgresRepo.UserRepo.FindOneByFilter(ctx, nil, filter)
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
			Id:          user.ID.String(),
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) CreateUser(ctx context.Context, data *oauth.CreateUserParams) (*common.SuccessResponse, error) {
	return nil, nil
}

// ------------------------ Helper ------------------------
func convertUserParamsToFilter(data *oauth.GetUserByFilterParams) (*repository.FindUserByFilter, error) {
	var userId uuid.UUID
	var userIds []uuid.UUID
	var err error
	limit := int(*data.Limit)
	offset := int(*data.Offset)

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
		ID:           &userId,
		Limit:        &limit,
		Offset:       &offset,
		IDs:          userIds,
		Emails:       data.Emails,
		PhoneNumbers: data.PhoneNumbers,
	}, nil
}
