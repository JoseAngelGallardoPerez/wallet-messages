package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetRolesResources checks roles and resources to be returned as expected
func TestServiceGetRolesResources(t *testing.T) {
	auth := createService()
	assert.IsType(t, map[string]map[string][]string{}, auth.getRolesResources())
}

// TestRegisterPermissions checks acl permissions works as expected
func TestServiceRegisterPermissions(t *testing.T) {
	auth := createService()
	auth.registerPermissions()

	for role, resources := range auth.getRolesResources() {
		for resource, actions := range resources {
			for _, action := range actions {
				assert.True(t, auth.Can(role, action, resource))
			}
		}
	}
}

// createService creates new service instance
func createService() *Service {
	return NewService(NewPermissionService())
}
