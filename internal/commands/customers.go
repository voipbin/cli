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

			items, nextToken, err := c.List(context.Background(), "/customers", params)
			if err != nil {
				return fmt.Errorf("could not list customers: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, customerListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/customers/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get customer: %w", err)
			}

			return output.PrintItem(cmd, result, customerDetailColumns)
		},
	}
}

func newCustomersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new customer",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
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

			body := map[string]interface{}{
				"name":           name,
				"email":          email,
				"phone_number":   phone,
				"address":        address,
				"detail":         detail,
				"webhook_uri":    webhookURI,
				"webhook_method": webhookMethod,
			}

			result, err := c.Post(context.Background(), "/customers", body)
			if err != nil {
				return fmt.Errorf("could not create customer: %w", err)
			}

			return output.PrintItem(cmd, result, customerDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
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

			body := map[string]interface{}{
				"name":           name,
				"email":          email,
				"phone_number":   phone,
				"address":        address,
				"detail":         detail,
				"webhook_uri":    webhookURI,
				"webhook_method": webhookMethod,
			}

			result, err := c.Put(context.Background(), "/customers/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update customer: %w", err)
			}

			return output.PrintItem(cmd, result, customerDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/customers/"+args[0]); err != nil {
				return fmt.Errorf("could not delete customer: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			billingAccountID, _ := cmd.Flags().GetString("billing-account-id")
			body := map[string]interface{}{
				"billing_account_id": billingAccountID,
			}

			result, err := c.Put(context.Background(), "/customers/"+args[0]+"/billing_account_id", body)
			if err != nil {
				return fmt.Errorf("could not update billing account: %w", err)
			}

			return output.PrintItem(cmd, result, customerDetailColumns)
		},
	}
	cmd.Flags().String("billing-account-id", "", "Billing account ID")
	_ = cmd.MarkFlagRequired("billing-account-id")
	return cmd
}
