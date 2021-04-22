package dao

func Providers() []interface{} {
	return []interface{}{
		NewMessage,
	}
}
