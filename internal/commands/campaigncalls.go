package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newCampaigncallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaigncalls",
		Short: "Manage campaign calls",
	}
	cmd.AddCommand(
		newCampaigncallsListCmd(),
		newCampaigncallsGetCmd(),
		newCampaigncallsDeleteCmd(),
	)
	return cmd
}

var campaigncallListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CAMPAIGN_ID", Field: "campaign_id"},
	{Name: "STATUS", Field: "status"},
	{Name: "RESULT", Field: "result"},
	{Name: "TRY_COUNT", Field: "try_count"},
	{Name: "CREATED", Field: "tm_create"},
}

var campaigncallDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "CAMPAIGN_ID", Field: "campaign_id"},
	{Name: "OUTDIAL_ID", Field: "outdial_id"},
	{Name: "OUTDIAL_TARGET_ID", Field: "outdial_target_id"},
	{Name: "OUTPLAN_ID", Field: "outplan_id"},
	{Name: "QUEUE_ID", Field: "queue_id"},
	{Name: "FLOW_ID", Field: "flow_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "RESULT", Field: "result"},
	{Name: "TRY_COUNT", Field: "try_count"},
	{Name: "CREATED", Field: "tm_create"},
}

func newCampaigncallsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List campaign calls",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetCampaigncallsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetCampaigncallsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list campaign calls: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, campaigncallListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newCampaigncallsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a campaign call by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetCampaigncallsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get campaign call: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, campaigncallDetailColumns)
		},
	}
}

func newCampaigncallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a campaign call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteCampaigncallsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete campaign call: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign call %s deleted.\n", args[0])
			return nil
		},
	}
}
