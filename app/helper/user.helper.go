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
) (*entity.User, error) {
	// Check user exited
	existedUser, err := h.postgresRepo.UserRepo.FindExistedByFilter(ctx, &repository.FindUserByFilter{
		PhoneNumber: data.PhoneNumber,
		Email:       data.Email,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUser) > 0 {
		return nil, errors.New(errors.ErrCodeUserExisted)
	}

	// Create user
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	user := entity.NewUser()
	user.PhoneNumber = data.PhoneNumber
	user.Email = data.Email
	user.Password = data.Password
	if err := h.postgresRepo.UserRepo.Create(ctx, tx, user); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return user, nil
}

func (h *userHelper) CreateUsers(
	ctx context.Context,
	data []CreateUserParams,
) ([]entity.User, error) {
	phoneNumbers := make([]string, 0)
	emails := make([]string, 0)
	users := make([]entity.User, 0)

	for _, item := range data {
		phoneNumbers = append(phoneNumbers, *item.PhoneNumber)
		emails = append(emails, *item.Email)
	}

	existedUser, err := h.postgresRepo.UserRepo.FindExistedByFilter(ctx, &repository.FindUserByFilter{
		PhoneNumbers: phoneNumbers,
		Emails:       emails,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if len(existedUser) > 0 {
		return nil, errors.New(errors.ErrCodeUserExisted)
	}

	// Create user
	tx := database.BeginPostgresTransaction()
	if tx.Error != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	for _, item := range data {
		user := entity.NewUser()
		user.PhoneNumber = item.PhoneNumber
		user.Email = item.Email
		user.Password = item.Password

		users = append(users, *user)
	}
	if err := h.postgresRepo.UserRepo.BulkCreate(ctx, tx, users); err != nil {
		tx.Rollback()
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return users, nil
}
