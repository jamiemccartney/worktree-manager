package consts

import (
	"os"
)

type ConfigDefaults struct {
	ConfigEditor            string
	AutomaticWorkOnAfterAdd bool
}

func GetConfigDefaults() ConfigDefaults {

	configEditor := os.Getenv("EDITOR")

	if configEditor == "" {
		configEditor = "vi"
	}

	return ConfigDefaults{
		ConfigEditor:            configEditor,
		AutomaticWorkOnAfterAdd: true,
	}
}
