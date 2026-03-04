package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newBillingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billings",
		Short: "View billings",
	}
	cmd.AddCommand(
		newBillingsListCmd(),
		newBillingsGetCmd(),
	)
	return cmd
}

var billingListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "ACCOUNT_ID", Field: "account_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "COST_TOTAL", Field: "cost_total"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var billingDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "ACCOUNT_ID", Field: "account_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "BILLING_UNIT_COUNT", Field: "billing_unit_count"},
	{Name: "COST_PER_UNIT", Field: "cost_per_unit"},
	{Name: "COST_TOTAL", Field: "cost_total"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

func newBillingsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List billings",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetBillingsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetBillingsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list billings: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, billingListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newBillingsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a billing by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			id, err := uuid.Parse(args[0])
			if err != nil {
				return fmt.Errorf("invalid UUID: %w", err)
			}

			resp, err := client.GetBillingsBillingIdWithResponse(context.Background(), id)
			if err != nil {
				return fmt.Errorf("could not get billing: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingDetailColumns)
		},
	}
}
