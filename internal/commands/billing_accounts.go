package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetBillingAccountsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetBillingAccountsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list billing accounts: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, billingAccountListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetBillingAccountsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get billing account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingAccountDetailColumns)
		},
	}
}

func newBillingAccountsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new billing account",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PostBillingAccountsJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}

			resp, err := client.PostBillingAccountsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create billing account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingAccountDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PutBillingAccountsIdJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}

			resp, err := client.PutBillingAccountsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update billing account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingAccountDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteBillingAccountsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete billing account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			amount, _ := cmd.Flags().GetFloat32("amount")
			body := voipbin_client.PostBillingAccountsIdBalanceAddForceJSONRequestBody{
				Balance: &amount,
			}

			resp, err := client.PostBillingAccountsIdBalanceAddForceWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not add balance: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingAccountDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			amount, _ := cmd.Flags().GetFloat32("amount")
			body := voipbin_client.PostBillingAccountsIdBalanceSubtractForceJSONRequestBody{
				Balance: &amount,
			}

			resp, err := client.PostBillingAccountsIdBalanceSubtractForceWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not subtract balance: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingAccountDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			paymentMethod, _ := cmd.Flags().GetString("payment-method")
			paymentType, _ := cmd.Flags().GetString("payment-type")

			body := voipbin_client.PutBillingAccountsIdPaymentInfoJSONRequestBody{}
			if paymentMethod != "" {
				pm := voipbin_client.BillingManagerAccountPaymentMethod(paymentMethod)
				body.PaymentMethod = &pm
			}
			if paymentType != "" {
				pt := voipbin_client.BillingManagerAccountPaymentType(paymentType)
				body.PaymentType = &pt
			}

			resp, err := client.PutBillingAccountsIdPaymentInfoWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update payment info: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, billingAccountDetailColumns)
		},
	}
	cmd.Flags().String("payment-method", "", "Payment method")
	cmd.Flags().String("payment-type", "", "Payment type")
	return cmd
}
