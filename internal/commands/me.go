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
	client, err := auth.NewClientFromContext(cmd)
	if err != nil {
		return err
	}

	resp, err := client.GetMeWithResponse(context.Background())
	if err != nil {
		return fmt.Errorf("could not get user info: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("API error: %s", resp.Status())
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("unexpected empty response")
	}

	columns := []output.Column{
		{Name: "ID", Field: "id"},
		{Name: "USERNAME", Field: "username"},
		{Name: "NAME", Field: "name"},
		{Name: "STATUS", Field: "status"},
		{Name: "PERMISSION", Field: "permission"},
		{Name: "CREATED", Field: "tm_create"},
	}

	return output.PrintItem(cmd, resp.JSON200, columns)
}
