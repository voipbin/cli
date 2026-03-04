package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetQueuecallsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetQueuecallsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list queuecalls: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil || resp.JSON200.Result == nil {
				return fmt.Errorf("unexpected empty response")
			}
			if resp.JSON200.NextPageToken != nil && *resp.JSON200.NextPageToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", *resp.JSON200.NextPageToken)
			}
			return output.PrintList(cmd, *resp.JSON200.Result, queuecallListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetQueuecallsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get queuecall: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, queuecallDetailColumns)
		},
	}
}

func newQueuecallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a queuecall",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteQueuecallsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete queuecall: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.PostQueuecallsIdKickWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not kick queuecall: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queuecall %s kicked.\n", args[0])
			return nil
		},
	}
}
