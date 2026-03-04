package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newActiveflowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activeflows",
		Short: "Manage active flows",
	}
	cmd.AddCommand(
		newActiveflowsListCmd(),
		newActiveflowsGetCmd(),
		newActiveflowsCreateCmd(),
		newActiveflowsDeleteCmd(),
		newActiveflowsStopCmd(),
	)
	return cmd
}

var activeflowListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "FLOW_ID", Field: "flow_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var activeflowDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "FLOW_ID", Field: "flow_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "FORWARD_ACTION_ID", Field: "forward_action_id"},
	{Name: "ON_COMPLETE_FLOW_ID", Field: "on_complete_flow_id"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newActiveflowsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List activeflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetActiveflowsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetActiveflowsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list activeflows: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, activeflowListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newActiveflowsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an activeflow by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetActiveflowsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get activeflow: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, activeflowDetailColumns)
		},
	}
}

func newActiveflowsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new activeflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			flowID, _ := cmd.Flags().GetString("flow-id")
			body := voipbin_client.PostActiveflowsJSONRequestBody{}
			if flowID != "" {
				body.FlowId = &flowID
			}
			resp, err := client.PostActiveflowsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create activeflow: %w", err)
			}
			if resp.StatusCode() != 201 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON201 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON201, activeflowDetailColumns)
		},
	}
	cmd.Flags().String("flow-id", "", "Flow ID to execute")
	_ = cmd.MarkFlagRequired("flow-id")
	return cmd
}

func newActiveflowsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an activeflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteActiveflowsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete activeflow: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Activeflow %s deleted.\n", args[0])
			return nil
		},
	}
}

func newActiveflowsStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <id>",
		Short: "Stop an activeflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.PostActiveflowsIdStopWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not stop activeflow: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Activeflow %s stopped.\n", args[0])
			return nil
		},
	}
}
