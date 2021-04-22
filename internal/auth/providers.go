package auth

func Providers() []interface{} {
	return []interface{}{
		NewService,
		NewPermissionService,
	}
}
