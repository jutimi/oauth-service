package service

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/package/database"
	"oauth-server/package/errors"
)

type userService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserService {

	return &userService{
		helpers:      helpers,
		postgresRepo: postgresRepo,
	}
}

func (s *userService) Register(ctx context.Context, data *model.RegisterRequest) (*model.RegisterResponse, error) {
	// Check user exited
	existedUser, err := s.postgresRepo.UserRepo.FindByFilter(ctx, nil, &repository.FindUserByFilter{
		PhoneNumber: &data.PhoneNumber,
		Email:       &data.Email,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUser) > 0 {
		return nil, errors.New(errors.ErrCodeUserExisted)
	}

	// Create user
	tx := database.BeginPostgresTransaction()
	user := entity.NewUser()
	user.PhoneNumber = &data.PhoneNumber
	user.Email = &data.Email
	user.Password = data.Password
	if err := s.postgresRepo.UserRepo.Create(ctx, tx, user); err != nil {
		tx.WithContext(ctx).Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	tx.WithContext(ctx).Commit()

	return &model.RegisterResponse{}, nil
}
