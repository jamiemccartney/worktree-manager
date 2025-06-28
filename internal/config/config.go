package config

import (
	"context"
	"fmt"
	"path/filepath"
	"worktree-manager/internal/consts"
	"worktree-manager/internal/fileops"
)

// Config represents the user configuration settings
type Config struct {
	ConfigEditor            string `json:"config-editor"`
	AutomaticWorkOnAfterAdd bool   `json:"automatic-work-on-after-add"`
}

var (
	cfg *Config
)

// Load reads and returns the configuration
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	var config Config
	if err := fileops.ReadJSONFile(consts.GetFilePaths().Config, &config); err != nil {
		return nil, err
	}

	cfg = &config
	return cfg, nil
}

// CheckConfigExists checks if the config file exists
func CheckConfigExists() bool {
	return fileops.FileExists(consts.GetFilePaths().Config)
}

// CreateDefault creates a default configuration file
func CreateDefault(force ...bool) error {
	shouldForce := len(force) > 0 && force[0]

	if CheckConfigExists() && !shouldForce {
		return fmt.Errorf("config file already exists at %s\n\n💡 Use 'wt init --force' to reinitialize", consts.GetFilePaths().Config)
	}

	defaults := consts.GetConfigDefaults()

	defaultConfig := Config{
		ConfigEditor:            defaults.ConfigEditor,
		AutomaticWorkOnAfterAdd: defaults.AutomaticWorkOnAfterAdd,
	}

	configPath := consts.GetFilePaths().Config
	if err := fileops.EnsureDir(filepath.Dir(configPath)); err != nil {
		return err
	}

	if err := fileops.WriteJSONFile(configPath, defaultConfig); err != nil {
		return err
	}

	cfg = &defaultConfig
	return nil
}

// Save saves the current configuration
func (c *Config) Save() error {
	return fileops.WriteJSONFile(consts.GetFilePaths().Config, c)
}

// GetConfigFromContext extracts config from context
func GetConfigFromContext(ctx context.Context) *Config {
	return ctx.Value(consts.GetContextKeys().Config).(*Config)
}
