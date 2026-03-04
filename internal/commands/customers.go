package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newCustomersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customers",
		Short: "Manage customers",
	}
	cmd.AddCommand(
		newCustomersListCmd(),
		newCustomersGetCmd(),
		newCustomersCreateCmd(),
		newCustomersUpdateCmd(),
		newCustomersDeleteCmd(),
		newCustomersUpdateBillingAccountIDCmd(),
	)
	return cmd
}

var customerListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "EMAIL", Field: "email"},
	{Name: "PHONE_NUMBER", Field: "phone_number"},
	{Name: "BILLING_ACCOUNT_ID", Field: "billing_account_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var customerDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "EMAIL", Field: "email"},
	{Name: "PHONE_NUMBER", Field: "phone_number"},
	{Name: "ADDRESS", Field: "address"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "BILLING_ACCOUNT_ID", Field: "billing_account_id"},
	{Name: "WEBHOOK_URI", Field: "webhook_uri"},
	{Name: "WEBHOOK_METHOD", Field: "webhook_method"},
	{Name: "CREATED", Field: "tm_create"},
}

func newCustomersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List customers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetCustomersParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetCustomersWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list customers: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, customerListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newCustomersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a customer by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetCustomersIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get customer: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, customerDetailColumns)
		},
	}
}

func newCustomersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new customer",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			phone, _ := cmd.Flags().GetString("phone-number")
			address, _ := cmd.Flags().GetString("address")
			detail, _ := cmd.Flags().GetString("detail")
			webhookURI, _ := cmd.Flags().GetString("webhook-uri")
			webhookMethod, _ := cmd.Flags().GetString("webhook-method")

			body := voipbin_client.PostCustomersJSONRequestBody{
				Name:          name,
				Email:         email,
				PhoneNumber:   phone,
				Address:       address,
				Detail:        detail,
				WebhookUri:    webhookURI,
				WebhookMethod: voipbin_client.CustomerManagerCustomerWebhookMethod(webhookMethod),
			}

			resp, err := client.PostCustomersWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create customer: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, customerDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Customer name")
	cmd.Flags().String("email", "", "Email address")
	cmd.Flags().String("phone-number", "", "Phone number")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("detail", "", "Detail")
	cmd.Flags().String("webhook-uri", "", "Webhook URI")
	cmd.Flags().String("webhook-method", "POST", "Webhook HTTP method")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newCustomersUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a customer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			email, _ := cmd.Flags().GetString("email")
			phone, _ := cmd.Flags().GetString("phone-number")
			address, _ := cmd.Flags().GetString("address")
			detail, _ := cmd.Flags().GetString("detail")
			webhookURI, _ := cmd.Flags().GetString("webhook-uri")
			webhookMethod, _ := cmd.Flags().GetString("webhook-method")

			body := voipbin_client.PutCustomersIdJSONRequestBody{
				Name:          name,
				Email:         email,
				PhoneNumber:   phone,
				Address:       address,
				Detail:        detail,
				WebhookUri:    webhookURI,
				WebhookMethod: voipbin_client.CustomerManagerCustomerWebhookMethod(webhookMethod),
			}

			resp, err := client.PutCustomersIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update customer: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, customerDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Customer name")
	cmd.Flags().String("email", "", "Email address")
	cmd.Flags().String("phone-number", "", "Phone number")
	cmd.Flags().String("address", "", "Address")
	cmd.Flags().String("detail", "", "Detail")
	cmd.Flags().String("webhook-uri", "", "Webhook URI")
	cmd.Flags().String("webhook-method", "", "Webhook HTTP method")
	return cmd
}

func newCustomersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a customer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteCustomersIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete customer: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Customer %s deleted.\n", args[0])
			return nil
		},
	}
}

func newCustomersUpdateBillingAccountIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-billing-account <id>",
		Short: "Update customer billing account ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			billingAccountID, _ := cmd.Flags().GetString("billing-account-id")
			body := voipbin_client.PutCustomersIdBillingAccountIdJSONRequestBody{
				BillingAccountId: billingAccountID,
			}

			resp, err := client.PutCustomersIdBillingAccountIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update billing account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, customerDetailColumns)
		},
	}
	cmd.Flags().String("billing-account-id", "", "Billing account ID")
	_ = cmd.MarkFlagRequired("billing-account-id")
	return cmd
}
