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

			items, nextToken, err := c.List(context.Background(), "/billings", params)
			if err != nil {
				return fmt.Errorf("could not list billings: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, billingListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/billings/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get billing: %w", err)
			}

			return output.PrintItem(cmd, result, billingDetailColumns)
		},
	}
}
