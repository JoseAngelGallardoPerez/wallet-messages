package auth

import (
	"log"

	userpb "github.com/Confialink/wallet-users/rpc/proto/users"
	goAcl "github.com/kildevaeld/go-acl"
)

const (
	ResourceMessages      = "private_messages"
	ResourceMessagesAdmin = "private_admin_messages"

	ActionCreate = "create"
	ActionUpdate = "update"
	ActionRead   = "read"
	ActionDelete = "delete"

	RoleRoot      = "root"
	RoleAdmin     = "admin"
	RoleClient    = "client"
)

type Service struct {
	Acl                *goAcl.ACL
	dynamicPermissions PermissionMap
	permissionService  *PermissionService
}

type Permission string

type Policy func(*userpb.User) bool
type PermissionMap map[string]map[string]map[string]Policy

func NewService(permissionService *PermissionService) *Service {
	a := goAcl.New(goAcl.NewMemoryStore())
	auth := Service{Acl: a, permissionService: permissionService}
	auth.registerPermissions()
	auth.dynamicPermissions = PermissionMap{
		RoleClient: {
			ResourceMessages: {ActionCreate: allowFunc},
		},
		RoleAdmin: {
			ResourceMessages: {ActionCreate: auth.ProvideCheckSpecificPermission(PermissionSendReplyMessage)},
		},
	}

	return &auth
}

// getRolesResources returns permissions data
func (auth *Service) getRolesResources() map[string]map[string][]string {
	return map[string]map[string][]string{
		RoleClient: {
			ResourceMessages: {ActionCreate, ActionUpdate, ActionRead, ActionDelete},
		},
		RoleAdmin: {
			ResourceMessages:      {ActionCreate, ActionUpdate, ActionRead, ActionDelete},
			ResourceMessagesAdmin: {ActionCreate, ActionUpdate, ActionRead, ActionDelete},
		},
		RoleRoot: {
			ResourceMessages:      {ActionCreate, ActionUpdate, ActionRead, ActionDelete},
			ResourceMessagesAdmin: {ActionCreate, ActionUpdate, ActionRead, ActionDelete},
		},
	}
}

// registerPermissions registers allowed actions for roles
func (auth *Service) registerPermissions() {
	for role, resources := range auth.getRolesResources() {
		auth.Acl.Role(role, "")
		for resource, actions := range resources {
			for _, action := range actions {
				auth.Acl.Allow(role, action, resource)
			}
		}
	}
}

// Can checks action is allowed
func (auth *Service) Can(role string, action string, resource string) bool {
	return auth.Acl.Can(role, action, resource)
}

// CanDynamic checks action is allowed by calling associated function
func (auth *Service) CanDynamic(user *userpb.User, action string, resourceName string) bool {
	if user.RoleName == RoleRoot {
		return true
	}

	function := auth.getPermissionFunc(user.RoleName, action, resourceName)
	return function(user)
}

// blockFunc always block access
func blockFunc(_ *userpb.User) bool {
	return false
}

// allowFunc always allows access
func allowFunc(_ *userpb.User) bool {
	return true
}

// getPermissionFunc returns function by role, action and resourceName.
// Returns blockFunc if proposed func not found
func (auth *Service) getPermissionFunc(role string, action string, resourceName string) Policy {
	if rolePermission, ok := auth.dynamicPermissions[role]; ok {
		if resourcePermission, ok := rolePermission[resourceName]; ok {
			if actionPermission, ok := resourcePermission[action]; ok {
				return actionPermission
			}
		}
	}
	return blockFunc
}

// CheckPermission calls permission service in order to check if user granted permission
func (auth *Service) CheckPermission(perm Permission, user *userpb.User) bool {
	result, err := auth.permissionService.Check(user.UID, string(perm))
	if err != nil {
		log.Printf("permission policy failed to check permission: %s", err.Error())
		return false
	}
	return result
}

func (auth *Service) ProvideCheckSpecificPermission(perm Permission) Policy {
	return func(user *userpb.User) bool {
		return auth.CheckPermission(perm, user)
	}
}
