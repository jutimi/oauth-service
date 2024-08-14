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
	"oauth-server/external/client"
	"oauth-server/package/errors"
	"oauth-server/utils"

	"github.com/jutimi/grpc-service/workspace"
)

type permissionService struct {
	helpers      helper.HelperCollections
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewPermissionService(
	helpers helper.HelperCollections,
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) PermissionService {
	return &permissionService{
		postgresRepo: postgresRepo,
		helpers:      helpers,
	}
}

func (s *permissionService) AddUserWorkspacePermission(
	ctx context.Context,
	data *model.AddUserWorkspacePermissionRequest,
) (*model.AddUserWorkspacePermissionResponse, error) {
	// Get user workspace data
	clientGRPC := client.NewWorkspaceClient()
	defer clientGRPC.CloseConn()

	userWorkspaceId := data.UserWorkspaceId.String()
	isActive := true
	userWorkspace, err := clientGRPC.GetUserWorkspaceByFilter(ctx, &workspace.GetUserWorkspaceByFilterParams{
		Id:       &userWorkspaceId,
		IsActive: &isActive,
	})
	if err != nil {
		return nil, err
	}
	userId, _ := utils.ConvertStringToUUID(userWorkspace.Data.UserId)
	WorkspaceId, _ := utils.ConvertStringToUUID(userWorkspace.Data.WorkspaceId)

	// Validate permission and get permission tree
	permissionMemo := make(map[string]bool, 0)
	for _, permission := range data.Permissions {
		if err := s.helpers.PermissionHelper.ValidatePermission(ctx, permission); err != nil {
			return nil, err
		}

		permissions := s.helpers.PermissionHelper.GetPermissions(ctx, permission)
		maps.Copy(permissionMemo, permissions)
	}
	permissionStr := ""
	for permission := range permissionMemo {
		permissionStr += fmt.Sprintf("%s/", permission)
	}

	// Check and remove old permission
	existedPermissions, err := s.postgresRepo.PermissionRepo.FindByFilter(ctx, &repository.FindPermissionByFilter{
		UserWorkspaceId: &data.UserWorkspaceId,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if s.postgresRepo.PermissionRepo.Delete(ctx, nil, &existedPermissions[0]) != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	// Create new permission
	userWorkspacePermission := entity.NewPermission()
	userWorkspacePermission.UserWorkspaceId = &data.UserWorkspaceId
	userWorkspacePermission.UserId = userId
	userWorkspacePermission.Scopes = permissionStr
	userWorkspacePermission.WorkspaceId = &WorkspaceId
	if err := s.postgresRepo.PermissionRepo.Create(ctx, nil, userWorkspacePermission); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.AddUserWorkspacePermissionResponse{}, nil
}

func (s *permissionService) RevokeUserWorkspacePermission(
	ctx context.Context,
	data *model.RevokeUserWorkspacePermissionRequest,
) (*model.RevokeUserWorkspacePermissionResponse, error) {
	existedPermissions, err := s.postgresRepo.PermissionRepo.FindOneByFilter(ctx, &repository.FindPermissionByFilter{
		UserWorkspaceId: &data.UserWorkspaceId,
	})
	if err != nil {
		return nil, errors.New(errors.ErrCodePermissionNotFound)
	}
	if s.postgresRepo.PermissionRepo.Delete(ctx, nil, existedPermissions) != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return &model.RevokeUserWorkspacePermissionResponse{}, nil
}

func (s *permissionService) GetPermissions(
	ctx context.Context,
) (*model.GetPermissionsResponse, error) {
	permissions := make([]model.PermissionDetail, 0)

	for permissionName, permission := range model.PERMISSION_TREE {
		for action := range permission {
			permissions = append(permissions, model.PermissionDetail{
				Name: fmt.Sprintf("%s_%s", action, permissionName),
				Key:  fmt.Sprintf("%s_%s", action, permissionName),
			})
		}
	}

	return &model.GetPermissionsResponse{
		Permissions: permissions,
	}, nil
}
