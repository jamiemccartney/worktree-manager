package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"worktree-manager/internal/output"

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
		homeDir = "."
	}
	configPath = filepath.Join(homeDir, ".worktree-manager", "config.json")
}

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

	config = resolveEnvVars(config)

	cfg = &config
	return cfg, nil
}

func GetConfigPath() string {
	return configPath
}

func ConfigExists() bool {
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

func CreateDefault(force ...bool) error {
	shouldForce := len(force) > 0 && force[0]

	if ConfigExists() && !shouldForce {
		return fmt.Errorf("config file already exists at %s\n\n💡 Use 'wt init --force' to reinitialize", configPath)
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

	configDir := filepath.Dir(configPath)

	output.Info("Worktree Directory created at: %s", configDir)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(defaultConfig, "", "    ")

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	output.Info("Config created at: %s", configPath)

	resolvedConfig := resolveEnvVars(defaultConfig)

	if err := os.MkdirAll(resolvedConfig.GitReposDir, 0755); err != nil {
		return fmt.Errorf("failed to create git repos directory: %w", err)
	}

	output.Info("Repos Directory created at: %s", resolvedConfig.GitReposDir)

	if err := os.MkdirAll(resolvedConfig.WorktreesDir, 0755); err != nil {
		return fmt.Errorf("failed to create worktrees directory: %w", err)
	}

	output.Info("Worktrees Directory created at: %s", resolvedConfig.WorktreesDir)
	output.Success("Init complete")

	return nil
}

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

func expandEnvVars(s string) string {
	if strings.Contains(s, "$HOME") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "."
		}
		s = strings.ReplaceAll(s, "$HOME", homeDir)
	}
	if strings.Contains(s, "$git-repos-dir") && cfg != nil {
		s = strings.ReplaceAll(s, "$git-repos-dir", cfg.GitReposDir)
	}
	return os.ExpandEnv(s)
}

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

func (c *Config) FindRepoByAlias(alias string) (*Repo, error) {
	for _, repo := range c.Repos {
		if repo.Alias == alias {
			return &repo, nil
		}
	}
	return nil, fmt.Errorf("repository with alias '%s' not found", alias)
}

func (c *Config) AddRepo(repo Repo) error {
	for _, existingRepo := range c.Repos {
		if existingRepo.Alias == repo.Alias {
			return fmt.Errorf("repository with alias '%s' already exists", repo.Alias)
		}
	}

	c.Repos = append(c.Repos, repo)
	return c.Save()
}

func (c *Config) RemoveRepo(alias string) error {
	for i, repo := range c.Repos {
		if repo.Alias == alias {
			c.Repos = append(c.Repos[:i], c.Repos[i+1:]...)
			return c.Save()
		}
	}
	return fmt.Errorf("repository with alias '%s' not found", alias)
}

func (c *Config) SetActiveRepo(alias string) error {
	_, err := c.FindRepoByAlias(alias)
	if err != nil {
		return err
	}

	c.ActiveRepo = alias
	return c.Save()
}

func (c *Config) GetActiveRepo() (*Repo, error) {
	if c.ActiveRepo == "" {
		return nil, fmt.Errorf("no active repository set. Use 'wt repo use <alias>' to set one")
	}

	return c.FindRepoByAlias(c.ActiveRepo)
}

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

func FormatRepoStatus(repo *Repo) string {
	status := ""

	if _, err := os.Stat(repo.Dir); os.IsNotExist(err) {
		status = "❌  Directory does not exist"
	} else {
		status = "✅  Available"

		worktreesDir := filepath.Join(repo.Dir, "worktrees")
		if entries, err := os.ReadDir(worktreesDir); err == nil {
			status += fmt.Sprintf(" (%d worktrees)", len(entries))
		}
	}

	return status
}

func PrintRepoList(repos []Repo) {
	if len(repos) == 0 {
		output.Hint("No repositories configured. Use 'wt repo clone <url>' to add one.")
		return
	}

	output.Info("Configured repositories (%d):", len(repos))

	for _, repo := range repos {
		output.Item(repo.Alias)
		output.Info("   Directory: %s", repo.Dir)
		output.Info("   Status: %s", FormatRepoStatus(&repo))

		if repo.PostWorktreeAddScript != "" {
			output.Info("   Post-add script: %s", repo.PostWorktreeAddScript)
		}
	}
}
