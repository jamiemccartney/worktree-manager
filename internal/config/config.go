package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ConfigEditor            string `json:"config-editor" mapstructure:"config-editor"`
	GitReposDir             string `json:"git-repos-dir" mapstructure:"git-repos-dir"`
	WorktreesDir            string `json:"worktrees-dir" mapstructure:"worktrees-dir"`
	AutomaticWorkOnAfterAdd bool   `json:"automatic-work-on-after-add" mapstructure:"automatic-work-on-after-add"`
	WorkOnScript            string `json:"work-on-script" mapstructure:"work-on-script"`
	ActiveRepo              string `json:"active-repo" mapstructure:"active-repo"`
	Repos                   []Repo `json:"repos" mapstructure:"repos"`
}

type Repo struct {
	Alias                 string `json:"alias" mapstructure:"alias"`
	Dir                   string `json:"dir" mapstructure:"dir"`
	PostWorktreeAddScript string `json:"post-worktree-add-script" mapstructure:"post-worktree-add-script"`
}

var (
	configPath string
	cfg        *Config
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to current directory if home directory cannot be determined
		homeDir = "."
	}
	configPath = filepath.Join(homeDir, ".worktree-manager", "config.json")
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Resolve environment variables
	config = resolveEnvVars(config)

	cfg = &config
	return cfg, nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	return configPath
}

// ConfigExists checks if the config file exists
func ConfigExists() bool {
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

// CreateDefault creates a default config file
func CreateDefault(force ...bool) error {
	shouldForce := len(force) > 0 && force[0]
	
	if ConfigExists() && !shouldForce {
		return fmt.Errorf("config file already exists at %s\n\n💡 Use 'wt init --force' to reinitialize", configPath)
	}
	
	if ConfigExists() && shouldForce {
		fmt.Printf("⚠️  Reinitializing existing configuration at %s\n", configPath)
	}

	defaultConfig := Config{
		ConfigEditor:            "vi",
		GitReposDir:             "$HOME/.worktree-manager/repos",
		WorktreesDir:            "$HOME/.worktree-manager/worktrees",
		AutomaticWorkOnAfterAdd: true,
		WorkOnScript:            "",
		ActiveRepo:              "",
		Repos:                   []Repo{},
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write config file
	data, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Create git repos and worktrees directories
	resolvedConfig := resolveEnvVars(defaultConfig)
	if err := os.MkdirAll(resolvedConfig.GitReposDir, 0755); err != nil {
		return fmt.Errorf("failed to create git repos directory: %w", err)
	}
	if err := os.MkdirAll(resolvedConfig.WorktreesDir, 0755); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	fmt.Printf("✅ Created default config at %s\n", configPath)
	fmt.Printf("✅ Created git repos directory at %s\n", resolvedConfig.GitReposDir)
	fmt.Printf("✅ Created worktrees directory at %s\n", resolvedConfig.WorktreesDir)
	return nil
}

// resolveEnvVars resolves environment variables in the config
func resolveEnvVars(config Config) Config {
	config.GitReposDir = expandEnvVars(config.GitReposDir)
	config.WorktreesDir = expandEnvVars(config.WorktreesDir)
	config.WorkOnScript = expandEnvVars(config.WorkOnScript)

	for i := range config.Repos {
		config.Repos[i].Dir = expandEnvVars(config.Repos[i].Dir)
		config.Repos[i].PostWorktreeAddScript = expandEnvVars(config.Repos[i].PostWorktreeAddScript)
	}

	return config
}

// expandEnvVars expands environment variables in a string
func expandEnvVars(s string) string {
	if strings.Contains(s, "$HOME") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fall back to current directory if home directory cannot be determined
			homeDir = "."
		}
		s = strings.ReplaceAll(s, "$HOME", homeDir)
	}
	if strings.Contains(s, "$git-repos-dir") && cfg != nil {
		s = strings.ReplaceAll(s, "$git-repos-dir", cfg.GitReposDir)
	}
	return os.ExpandEnv(s)
}

// GetCurrentRepo returns the current repository based on the current working directory
func (c *Config) GetCurrentRepo() (*Repo, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	for _, repo := range c.Repos {
		if strings.HasPrefix(pwd, repo.Dir) {
			return &repo, nil
		}
	}

	return nil, fmt.Errorf("current directory is not within a managed repository")
}

// FindRepoByAlias finds a repository by its alias
func (c *Config) FindRepoByAlias(alias string) (*Repo, error) {
	for _, repo := range c.Repos {
		if repo.Alias == alias {
			return &repo, nil
		}
	}
	return nil, fmt.Errorf("repository with alias '%s' not found", alias)
}

// AddRepo adds a new repository to the config
func (c *Config) AddRepo(repo Repo) error {
	// Check if alias already exists
	for _, existingRepo := range c.Repos {
		if existingRepo.Alias == repo.Alias {
			return fmt.Errorf("repository with alias '%s' already exists", repo.Alias)
		}
	}

	c.Repos = append(c.Repos, repo)
	return c.Save()
}

// RemoveRepo removes a repository from the config
func (c *Config) RemoveRepo(alias string) error {
	for i, repo := range c.Repos {
		if repo.Alias == alias {
			c.Repos = append(c.Repos[:i], c.Repos[i+1:]...)
			return c.Save()
		}
	}
	return fmt.Errorf("repository with alias '%s' not found", alias)
}

// SetActiveRepo sets the active repository by alias
func (c *Config) SetActiveRepo(alias string) error {
	// Verify the repository exists
	_, err := c.FindRepoByAlias(alias)
	if err != nil {
		return err
	}
	
	c.ActiveRepo = alias
	return c.Save()
}

// GetActiveRepo returns the active repository
func (c *Config) GetActiveRepo() (*Repo, error) {
	if c.ActiveRepo == "" {
		return nil, fmt.Errorf("no active repository set. Use 'wt repo use <alias>' to set one")
	}
	
	return c.FindRepoByAlias(c.ActiveRepo)
}

// Save saves the current config to the file
func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
