package service

func Providers() []interface{} {
	return []interface{}{
		NewUserService,
		NewNotificationService,
		NewMessage,
		NewCsv,
	}
}
