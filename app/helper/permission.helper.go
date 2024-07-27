package helper

import (
	"oauth-server/app/model"
	postgres_repository "oauth-server/app/repository/postgres"
	"oauth-server/package/errors"
	"strings"
)

type permissionHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewPermissionHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) PermissionHelper {
	return &permissionHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *permissionHelper) ValidatePermission(permission string) error {
	permissionArr := strings.Split(permission, "_")
	action := permissionArr[0]
	resource := strings.Join(permissionArr[1:], "_")

	permissionDetail, ok := model.PERMISSION_TREE[resource]
	if !ok {
		return errors.New(errors.ErrCodePermissionNotFound)
	}

	if _, ok := permissionDetail[action]; !ok {
		return errors.New(errors.ErrCodePermissionActionNotFound)
	}

	return nil
}

func (h *permissionHelper) GetPermissions(permission string) map[string]bool {
	result := make(map[string]bool)

	stacks := []string{permission}

	for len(stacks) > 0 {
		current := stacks[len(stacks)-1]
		stacks = stacks[:len(stacks)-1]

		permissionArr := strings.Split(current, "_")
		action := permissionArr[0]
		resource := strings.Join(permissionArr[1:], "_")

		// Check if the permission is already in the result
		if _, ok := result[current]; ok {
			continue
		}

		// Add the permission to the result
		result[current] = true

		// Get the next permissions
		nextPermissions, exists := model.PERMISSION_TREE[resource][action]
		if !exists {
			continue
		}

		// Add the next permissions to the stack
		for _, nextPermission := range nextPermissions {
			stacks = append(stacks, nextPermission)
		}
	}

	return result
}
