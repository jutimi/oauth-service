package model

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	USER_WORKSPACE_PERMISSION = "user_workspaces"
	ORGANIZATION_PERMISSION   = "organizations"
	SHIFT_PERMISSION          = "shifts"
)

const (
	PERMISSION_READ   = "get"
	PERMISSION_CREATE = "create"
	PERMISSION_UPDATE = "update"
	PERMISSION_DELETE = "delete"
	PERMISSION_BYPASS = "bypass"
)

var PERMISSION_RESOURCES = []string{
	USER_WORKSPACE_PERMISSION,
	ORGANIZATION_PERMISSION,
	SHIFT_PERMISSION,
}

var PERMISSION_TREE = map[string]map[string][]string{
	USER_WORKSPACE_PERMISSION: {
		PERMISSION_READ: nil,
		PERMISSION_CREATE: {
			fmt.Sprintf("%s_%s", PERMISSION_READ, USER_WORKSPACE_PERMISSION),
		},
		PERMISSION_UPDATE: {
			fmt.Sprintf("%s_%s", PERMISSION_CREATE, USER_WORKSPACE_PERMISSION),
		},
		PERMISSION_DELETE: {
			fmt.Sprintf("%s_%s", PERMISSION_UPDATE, USER_WORKSPACE_PERMISSION),
		},
	},
	ORGANIZATION_PERMISSION: {
		PERMISSION_READ: {
			fmt.Sprintf("%s_%s", PERMISSION_READ, USER_WORKSPACE_PERMISSION),
		},
		PERMISSION_CREATE: {
			fmt.Sprintf("%s_%s", PERMISSION_READ, ORGANIZATION_PERMISSION),
		},
		PERMISSION_UPDATE: {
			fmt.Sprintf("%s_%s", PERMISSION_CREATE, ORGANIZATION_PERMISSION),
		},
		PERMISSION_DELETE: {
			fmt.Sprintf("%s_%s", PERMISSION_UPDATE, ORGANIZATION_PERMISSION),
		},
	},
}

type AddUserWorkspacePermissionRequest struct {
	UserWorkspaceId uuid.UUID `json:"user_workspace_id" validate:"required,uuid"`
	Permissions     []string  `json:"permissions" validate:"required"`
}
type AddUserWorkspacePermissionResponse struct{}

type RevokeUserWorkspacePermissionRequest struct {
	UserWorkspaceId uuid.UUID `json:"user_workspace_id" validate:"required,uuid"`
}
type RevokeUserWorkspacePermissionResponse struct{}

type GetPermissionsResponse struct {
	Permissions []PermissionDetail `json:"permissions"`
}
type PermissionDetail struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
