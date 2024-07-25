package helper

import (
	"context"
	"oauth-server/app/entity"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/package/database"
	"oauth-server/package/errors"
)

type userHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserHelper {
	return &userHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *userHelper) CreateUser(
	ctx context.Context,
	data *CreateUserParams,
) error {
	// Check user exited
	existedUser, err := h.postgresRepo.UserRepo.FindExistedByFilter(ctx, nil, &repository.FindUserByFilter{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	})
	if err != nil {
		return errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUser) > 0 {
		return errors.New(errors.ErrCodeUserExisted)
	}

	// Create user
	tx := database.BeginPostgresTransaction()
	user := entity.NewUser()
	user.PhoneNumber = data.PhoneNumber
	user.Email = data.Email
	user.Password = data.Password
	if err := h.postgresRepo.UserRepo.Create(ctx, tx, user); err != nil {
		tx.WithContext(ctx).Rollback()
		return errors.New(errors.ErrCodeInternalServerError)
	}
	tx.WithContext(ctx).Commit()

	return nil
}
