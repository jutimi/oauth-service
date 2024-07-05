package user_grpc

import (
	context "context"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/utils"

	"github.com/google/uuid"
	"google.golang.org/grpc/status"
)

type userServer struct {
	UnimplementedUserRouteServer

	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserServer(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserRouteServer {
	return &userServer{
		postgresRepo: postgresRepo,
	}
}

func (s *userServer) GetUserById(ctx context.Context, data *GetUserByIdParams) (*UserResponse, error) {
	userId, err := utils.ConvertStringToUUID(data.Id)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}

	user, err := s.postgresRepo.PostgresUserRepo.FindUserByFilter(ctx, nil, &repository.FindUserByFilter{
		ID: &userId,
	})
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}

	return &UserResponse{
		Id:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}, nil
}

func (s *userServer) GetUsersByFilter(ctx context.Context, data *GetUserByFilterParams) (*UsersResponse, error) {

	var usersRes []*UserResponse

	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}

	users, err := s.postgresRepo.PostgresUserRepo.FindUsersByFilter(ctx, nil, filter)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}

	for _, user := range users {
		usersRes = append(usersRes, &UserResponse{
			Id:          user.ID.String(),
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		})
	}

	return &UsersResponse{
		Users: usersRes,
	}, nil
}

func (s *userServer) GetUserByFilter(ctx context.Context, data *GetUserByFilterParams) (*UserResponse, error) {
	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}

	user, err := s.postgresRepo.PostgresUserRepo.FindUserByFilter(ctx, nil, filter)
	if err != nil {
		return nil, status.Error(status.Code(err), err.Error())
	}

	return &UserResponse{
		Id:          user.ID.String(),
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
	}, nil
}

// ------------------------ Helper ------------------------
func convertUserParamsToFilter(data *GetUserByFilterParams) (*repository.FindUserByFilter, error) {
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
