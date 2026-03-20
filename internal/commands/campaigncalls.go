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

			items, nextToken, err := c.List(context.Background(), "/campaigncalls", params)
			if err != nil {
				return fmt.Errorf("could not list campaign calls: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, campaigncallListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/campaigncalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get campaign call: %w", err)
			}

			return output.PrintItem(cmd, result, campaigncallDetailColumns)
		},
	}
}

func newCampaigncallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a campaign call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/campaigncalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete campaign call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign call %s deleted.\n", args[0])
			return nil
		},
	}
}
