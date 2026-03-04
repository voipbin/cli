package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultAPIURL = "https://api.voipbin.net/v1.0"
	configDir     = ".vn"
	configFile    = "config.yaml"
)

type Profile struct {
	AccessKey string `yaml:"access_key"`
	APIURL    string `yaml:"api_url,omitempty"`
}

type Config struct {
	CurrentProfile string             `yaml:"current_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, configDir, configFile), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				CurrentProfile: "default",
				Profiles:       make(map[string]Profile),
			}, nil
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not write config: %w", err)
	}
	return nil
}

func (c *Config) GetProfile(name string) (Profile, bool) {
	p, ok := c.Profiles[name]
	return p, ok
}

func (c *Config) SetProfile(name string, profile Profile) {
	c.Profiles[name] = profile
}

func (c *Config) DeleteProfile(name string) {
	delete(c.Profiles, name)
	if c.CurrentProfile == name {
		c.CurrentProfile = "default"
	}
}

func (c *Config) CurrentProfileData() (Profile, bool) {
	return c.GetProfile(c.CurrentProfile)
}
