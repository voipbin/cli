package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNonExistent(t *testing.T) {
	// Override home to temp dir
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if cfg.CurrentProfile != "default" {
		t.Errorf("expected default profile, got %q", cfg.CurrentProfile)
	}
	if len(cfg.Profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(cfg.Profiles))
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cfg := &Config{
		CurrentProfile: "test",
		Profiles: map[string]Profile{
			"test": {AccessKey: "ak-123", APIURL: "https://test.api.com"},
		},
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Verify file permissions
	path := filepath.Join(tmp, ".vn", "config.yaml")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected 0600 perms, got %o", info.Mode().Perm())
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if loaded.CurrentProfile != "test" {
		t.Errorf("expected profile 'test', got %q", loaded.CurrentProfile)
	}
	p, ok := loaded.GetProfile("test")
	if !ok {
		t.Fatal("profile 'test' not found")
	}
	if p.AccessKey != "ak-123" {
		t.Errorf("expected access key 'ak-123', got %q", p.AccessKey)
	}
	if p.APIURL != "https://test.api.com" {
		t.Errorf("expected API URL 'https://test.api.com', got %q", p.APIURL)
	}
}

func TestProfileOperations(t *testing.T) {
	cfg := &Config{
		CurrentProfile: "default",
		Profiles:       map[string]Profile{},
	}

	cfg.SetProfile("prod", Profile{AccessKey: "ak-prod"})
	cfg.SetProfile("staging", Profile{AccessKey: "ak-staging"})

	if _, ok := cfg.GetProfile("prod"); !ok {
		t.Error("prod profile not found")
	}
	if _, ok := cfg.GetProfile("staging"); !ok {
		t.Error("staging profile not found")
	}
	if _, ok := cfg.GetProfile("nonexistent"); ok {
		t.Error("nonexistent profile should not exist")
	}

	cfg.CurrentProfile = "prod"
	p, ok := cfg.CurrentProfileData()
	if !ok {
		t.Fatal("current profile not found")
	}
	if p.AccessKey != "ak-prod" {
		t.Errorf("expected ak-prod, got %q", p.AccessKey)
	}

	cfg.DeleteProfile("prod")
	if _, ok := cfg.GetProfile("prod"); ok {
		t.Error("prod profile should be deleted")
	}
	if cfg.CurrentProfile != "default" {
		t.Errorf("current profile should reset to default after deleting active profile, got %q", cfg.CurrentProfile)
	}
}
