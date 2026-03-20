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

func newBillingAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "billing-accounts",
		Short: "Manage billing accounts",
	}
	cmd.AddCommand(
		newBillingAccountsListCmd(),
		newBillingAccountsGetCmd(),
		newBillingAccountsCreateCmd(),
		newBillingAccountsUpdateCmd(),
		newBillingAccountsDeleteCmd(),
		newBillingAccountsBalanceAddCmd(),
		newBillingAccountsBalanceSubtractCmd(),
		newBillingAccountsUpdatePaymentInfoCmd(),
	)
	return cmd
}

var billingAccountListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "BALANCE", Field: "balance"},
	{Name: "PAYMENT_TYPE", Field: "payment_type"},
	{Name: "CREATED", Field: "tm_create"},
}

var billingAccountDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "BALANCE", Field: "balance"},
	{Name: "PAYMENT_TYPE", Field: "payment_type"},
	{Name: "PAYMENT_METHOD", Field: "payment_method"},
	{Name: "CREATED", Field: "tm_create"},
}

func newBillingAccountsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List billing accounts",
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

			items, nextToken, err := c.List(context.Background(), "/billing_accounts", params)
			if err != nil {
				return fmt.Errorf("could not list billing accounts: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, billingAccountListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newBillingAccountsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a billing account by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/billing_accounts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get billing account: %w", err)
			}

			return output.PrintItem(cmd, result, billingAccountDetailColumns)
		},
	}
}

func newBillingAccountsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new billing account",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := map[string]interface{}{
				"name":   name,
				"detail": detail,
			}

			result, err := c.Post(context.Background(), "/billing_accounts", body)
			if err != nil {
				return fmt.Errorf("could not create billing account: %w", err)
			}

			return output.PrintItem(cmd, result, billingAccountDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Account name")
	cmd.Flags().String("detail", "", "Account detail")
	return cmd
}

func newBillingAccountsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a billing account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}

			result, err := c.Put(context.Background(), "/billing_accounts/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update billing account: %w", err)
			}

			return output.PrintItem(cmd, result, billingAccountDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Account name")
	cmd.Flags().String("detail", "", "Account detail")
	return cmd
}

func newBillingAccountsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a billing account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/billing_accounts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete billing account: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Billing account %s deleted.\n", args[0])
			return nil
		},
	}
}

func newBillingAccountsBalanceAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance-add <id>",
		Short: "Add balance to a billing account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			amount, _ := cmd.Flags().GetFloat32("amount")
			body := map[string]interface{}{
				"balance": amount,
			}

			result, err := c.Post(context.Background(), "/billing_accounts/"+args[0]+"/balance_add_force", body)
			if err != nil {
				return fmt.Errorf("could not add balance: %w", err)
			}

			return output.PrintItem(cmd, result, billingAccountDetailColumns)
		},
	}
	cmd.Flags().Float32("amount", 0, "Amount to add (USD)")
	_ = cmd.MarkFlagRequired("amount")
	return cmd
}

func newBillingAccountsBalanceSubtractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance-subtract <id>",
		Short: "Subtract balance from a billing account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			amount, _ := cmd.Flags().GetFloat32("amount")
			body := map[string]interface{}{
				"balance": amount,
			}

			result, err := c.Post(context.Background(), "/billing_accounts/"+args[0]+"/balance_subtract_force", body)
			if err != nil {
				return fmt.Errorf("could not subtract balance: %w", err)
			}

			return output.PrintItem(cmd, result, billingAccountDetailColumns)
		},
	}
	cmd.Flags().Float32("amount", 0, "Amount to subtract (USD)")
	_ = cmd.MarkFlagRequired("amount")
	return cmd
}

func newBillingAccountsUpdatePaymentInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-payment-info <id>",
		Short: "Update payment info for a billing account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			paymentMethod, _ := cmd.Flags().GetString("payment-method")
			paymentType, _ := cmd.Flags().GetString("payment-type")

			body := map[string]interface{}{}
			if paymentMethod != "" {
				body["payment_method"] = paymentMethod
			}
			if paymentType != "" {
				body["payment_type"] = paymentType
			}

			result, err := c.Put(context.Background(), "/billing_accounts/"+args[0]+"/payment_info", body)
			if err != nil {
				return fmt.Errorf("could not update payment info: %w", err)
			}

			return output.PrintItem(cmd, result, billingAccountDetailColumns)
		},
	}
	cmd.Flags().String("payment-method", "", "Payment method")
	cmd.Flags().String("payment-type", "", "Payment type")
	return cmd
}
