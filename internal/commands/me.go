package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
)

func newMeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "me",
		Short: "Show current authenticated user",
		RunE:  runMe,
	}
}

func runMe(cmd *cobra.Command, args []string) error {
	c, err := auth.NewClientFromContext(cmd)
	if err != nil {
		return err
	}

	result, err := c.Get(context.Background(), "/me")
	if err != nil {
		return fmt.Errorf("could not get user info: %w", err)
	}

	columns := []output.Column{
		{Name: "ID", Field: "id"},
		{Name: "USERNAME", Field: "username"},
		{Name: "NAME", Field: "name"},
		{Name: "STATUS", Field: "status"},
		{Name: "PERMISSION", Field: "permission"},
		{Name: "CREATED", Field: "tm_create"},
	}

	return output.PrintItem(cmd, result, columns)
}
