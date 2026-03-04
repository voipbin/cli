package commands

import (
	"strings"
	"testing"
)

func TestRootCommandHasAllResources(t *testing.T) {
	root := NewRootCmd()

	expectedCommands := []string{
		"accesskeys", "activeflows", "agents", "aicalls", "aimessages",
		"ais", "aisummaries", "available-numbers", "billing-accounts",
		"billings", "calls", "campaigncalls", "campaigns", "chatmessages",
		"chatroommessages", "chatrooms", "chats", "completion",
		"conferencecalls", "conferences", "conversation-accounts",
		"conversations", "customers", "emails", "extensions", "files",
		"flows", "groupcalls", "login", "logout", "me", "messages",
		"numbers", "outdials", "outplans", "providers", "queuecalls",
		"queues", "recordings", "routes", "storage-accounts",
		"storage-files", "tags", "transcribes", "transfers", "trunks",
		"version",
	}

	registered := make(map[string]bool)
	for _, cmd := range root.Commands() {
		registered[cmd.Name()] = true
	}

	for _, name := range expectedCommands {
		if !registered[name] {
			t.Errorf("command %q not registered in root", name)
		}
	}
}

func TestResourceCommandsHaveSubcommands(t *testing.T) {
	root := NewRootCmd()

	resources := []string{
		"calls", "messages", "emails", "conferences", "agents",
		"campaigns", "flows", "queues", "numbers", "recordings",
		"tags", "chats", "ais", "trunks", "providers",
	}

	for _, name := range resources {
		cmd, _, err := root.Find([]string{name})
		if err != nil {
			t.Errorf("command %q not found: %v", name, err)
			continue
		}
		if !cmd.HasSubCommands() {
			t.Errorf("command %q has no subcommands", name)
		}
	}
}

func TestGlobalFlags(t *testing.T) {
	root := NewRootCmd()

	flags := []string{"output", "profile", "access-key", "api-url"}
	for _, name := range flags {
		f := root.PersistentFlags().Lookup(name)
		if f == nil {
			t.Errorf("global flag %q not found", name)
		}
	}
}

func TestCallsSubcommands(t *testing.T) {
	root := NewRootCmd()
	callsCmd, _, _ := root.Find([]string{"calls"})

	expected := []string{
		"list", "get", "create", "delete", "hangup",
		"hold", "unhold", "mute", "unmute",
	}

	registered := make(map[string]bool)
	for _, cmd := range callsCmd.Commands() {
		registered[cmd.Name()] = true
	}

	for _, name := range expected {
		if !registered[name] {
			t.Errorf("calls subcommand %q not registered", name)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	root := NewRootCmd()
	out := &strings.Builder{}
	root.SetOut(out)
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if !strings.Contains(out.String(), "vn version") {
		t.Errorf("version output missing expected text, got: %s", out.String())
	}
}
