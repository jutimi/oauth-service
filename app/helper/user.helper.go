package helper

import (
	"context"
	"oauth-server/app/entity"
)

type userHelper struct {
}

func NewUserHelper() UserHelper {
	return &userHelper{}
}

func (h *userHelper) CreateUser(
	ctx context.Context,
	user *entity.User,
) error {
	return nil
}
