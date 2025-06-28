package consts

type ContextKey string

type ContextKeyNames struct {
	Config ContextKey
	State  ContextKey
}

func GetContextKeys() ContextKeyNames {
	return ContextKeyNames{
		Config: "config",
		State:  "state",
	}
}
