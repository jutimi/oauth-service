package service

import (
	"context"
	"fmt"
	"maps"
	"oauth-server/app/entity"
	"oauth-server/app/helper"
	"oauth-server/app/model"
	"oauth-server/app/repository"
	postgres_repository "oauth-server/app/repository/postgres"
	client_grpc "oauth-server/grpc/client"
	"oauth-server/package/errors"
	"oauth-server/utils"

	"github.com/jutimi/grpc-service/workspace"
)

type permissionService struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
	helpers      helper.HelperCollections
}

func NewPermissionService(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	helpers helper.HelperCollections,
) PermissionService {
	return &permissionService{
		postgresRepo: postgresRepo,
		helpers:      helpers,
	}
}

func (s *permissionService) AddUserWSPermission(
	ctx context.Context,
	data *model.AddUserWSPermissionRequest,
) (*model.AddUserWSPermissionResponse, error) {
	// Get user workspace data
	clientGRPC := client_grpc.NewWsClient()
	defer clientGRPC.CloseConn()

	userWSId := data.UserWorkspaceId.String()
	isActive := true
	userWS, err := clientGRPC.GetUserWSByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
		Id:       &userWSId,
		IsActive: &isActive,
	})
	if err != nil {
		return nil, err
	}
	userId, _ := utils.ConvertStringToUUID(userWS.Data.UserId)
	wsId, _ := utils.ConvertStringToUUID(userWS.Data.WorkspaceId)

	// Validate permission and get permission tree
	permissionMemo := make(map[string]bool, 0)
	for _, permission := range data.Permissions {
		if err := s.helpers.PermissionHelper.ValidatePermission(permission); err != nil {
			return nil, err
		}

		permissions := s.helpers.PermissionHelper.GetPermissions(permission)
		maps.Copy(permissionMemo, permissions)
	}
	permissionStr := ""
	for permission := range permissionMemo {
		permissionStr += fmt.Sprintf("%s/", permission)
	}

	// Check and remove old permission
	existedPermissions, err := s.postgresRepo.PermissionRepo.FindByFilter(ctx, &repository.FindPermissionByFilter{
		UserWorkspaceID: &data.UserWorkspaceId,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if s.postgresRepo.PermissionRepo.Delete(ctx, nil, &existedPermissions[0]) != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Create new permission
	userWSPermission := entity.NewPermission()
	userWSPermission.UserWorkspaceID = &data.UserWorkspaceId
	userWSPermission.UserID = userId
	userWSPermission.Scopes = permissionStr
	userWSPermission.WorkspaceID = &wsId
	if err := s.postgresRepo.PermissionRepo.Create(ctx, nil, userWSPermission); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.AddUserWSPermissionResponse{}, nil
}

func (s *permissionService) RevokeUserWSPermission(
	ctx context.Context,
	data *model.RevokeUserWSPermissionRequest,
) (*model.RevokeUserWSPermissionResponse, error) {
	existedPermissions, err := s.postgresRepo.PermissionRepo.FindOneByFilter(ctx, &repository.FindPermissionByFilter{
		UserWorkspaceID: &data.UserWorkspaceId,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodePermissionNotFound)
	}
	if s.postgresRepo.PermissionRepo.Delete(ctx, nil, existedPermissions) != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.RevokeUserWSPermissionResponse{}, nil
}
