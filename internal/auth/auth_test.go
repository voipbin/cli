package auth

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/config"
)

// newTestCmd creates a command with flags matching the real root command.
// We use Flags() (not PersistentFlags) so they're directly accessible
// without cobra's parent-child flag merging that happens during Execute().
func newTestCmd() *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Flags().String("access-key", "", "")
	cmd.Flags().String("profile", "", "")
	cmd.Flags().String("api-url", "", "")
	return cmd
}

func TestResolveAccessKeyFromFlag(t *testing.T) {
	cmd := newTestCmd()
	_ = cmd.Flags().Set("access-key", "ak-from-flag")
	t.Setenv("HOME", t.TempDir())

	key, err := resolveAccessKey(cmd)
	if err != nil {
		t.Fatalf("resolveAccessKey error: %v", err)
	}
	if key != "ak-from-flag" {
		t.Errorf("expected ak-from-flag, got %q", key)
	}
}

func TestResolveAccessKeyFromEnv(t *testing.T) {
	t.Setenv("VN_ACCESS_KEY", "ak-from-env")
	t.Setenv("HOME", t.TempDir())

	cmd := newTestCmd()

	key, err := resolveAccessKey(cmd)
	if err != nil {
		t.Fatalf("resolveAccessKey error: %v", err)
	}
	if key != "ak-from-env" {
		t.Errorf("expected ak-from-env, got %q", key)
	}
}

func TestResolveAccessKeyFromConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("VN_ACCESS_KEY", "")

	cfg := &config.Config{
		CurrentProfile: "default",
		Profiles: map[string]config.Profile{
			"default": {AccessKey: "ak-from-config"},
		},
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save config error: %v", err)
	}

	cmd := newTestCmd()

	key, err := resolveAccessKey(cmd)
	if err != nil {
		t.Fatalf("resolveAccessKey error: %v", err)
	}
	if key != "ak-from-config" {
		t.Errorf("expected ak-from-config, got %q", key)
	}
}

func TestResolveAccessKeyPriority(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("VN_ACCESS_KEY", "ak-from-env")

	cfg := &config.Config{
		CurrentProfile: "default",
		Profiles: map[string]config.Profile{
			"default": {AccessKey: "ak-from-config"},
		},
	}
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save config error: %v", err)
	}

	// Flag > env > config
	cmd := newTestCmd()
	_ = cmd.Flags().Set("access-key", "ak-from-flag")

	key, err := resolveAccessKey(cmd)
	if err != nil {
		t.Fatalf("resolveAccessKey error: %v", err)
	}
	if key != "ak-from-flag" {
		t.Errorf("flag should win: expected ak-from-flag, got %q", key)
	}
}

func TestResolveAccessKeyMissing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("VN_ACCESS_KEY", "")

	cmd := newTestCmd()

	_, err := resolveAccessKey(cmd)
	if err == nil {
		t.Error("expected error when no access key available")
	}
}

func TestResolveAPIURL(t *testing.T) {
	cmd := newTestCmd()
	t.Setenv("HOME", t.TempDir())

	url := resolveAPIURL(cmd)
	if url != config.DefaultAPIURL {
		t.Errorf("expected default URL, got %q", url)
	}

	_ = cmd.Flags().Set("api-url", "https://custom.api.com")
	url = resolveAPIURL(cmd)
	if url != "https://custom.api.com" {
		t.Errorf("expected custom URL, got %q", url)
	}
}
