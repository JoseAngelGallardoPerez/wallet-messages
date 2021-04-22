package rpc

func Providers() []interface{} {
	return []interface{}{
		NewPbServer,
	}
}
