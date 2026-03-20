package commands

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
)

func newQueuecallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queuecalls",
		Short: "Manage queue calls",
	}
	cmd.AddCommand(
		newQueuecallsListCmd(),
		newQueuecallsGetCmd(),
		newQueuecallsDeleteCmd(),
		newQueuecallsKickCmd(),
	)
	return cmd
}

var queuecallListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "SERVICE_AGENT_ID", Field: "service_agent_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var queuecallDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "SERVICE_AGENT_ID", Field: "service_agent_id"},
	{Name: "DURATION_WAITING", Field: "duration_waiting"},
	{Name: "DURATION_SERVICE", Field: "duration_service"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newQueuecallsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List queuecalls",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := url.Values{}
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}
			items, nextToken, err := c.List(context.Background(), "/queuecalls", params)
			if err != nil {
				return fmt.Errorf("could not list queuecalls: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, queuecallListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newQueuecallsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a queuecall by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/queuecalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get queuecall: %w", err)
			}
			return output.PrintItem(cmd, item, queuecallDetailColumns)
		},
	}
}

func newQueuecallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a queuecall",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/queuecalls/"+args[0]); err != nil {
				return fmt.Errorf("could not delete queuecall: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queuecall %s deleted.\n", args[0])
			return nil
		},
	}
}

func newQueuecallsKickCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kick <id>",
		Short: "Kick a queuecall",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Post(context.Background(), "/queuecalls/"+args[0]+"/kick", nil); err != nil {
				return fmt.Errorf("could not kick queuecall: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queuecall %s kicked.\n", args[0])
			return nil
		},
	}
}
