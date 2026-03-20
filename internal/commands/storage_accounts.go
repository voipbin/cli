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

func newStorageAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage-accounts",
		Short: "Manage storage accounts",
	}
	cmd.AddCommand(
		newStorageAccountsListCmd(),
		newStorageAccountsGetCmd(),
		newStorageAccountsCreateCmd(),
		newStorageAccountsDeleteCmd(),
	)
	return cmd
}

var storageAccountListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "TOTAL_FILE_COUNT", Field: "total_file_count"},
	{Name: "TOTAL_FILE_SIZE", Field: "total_file_size"},
	{Name: "CREATED", Field: "tm_create"},
}

var storageAccountDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "TOTAL_FILE_COUNT", Field: "total_file_count"},
	{Name: "TOTAL_FILE_SIZE", Field: "total_file_size"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newStorageAccountsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List storage accounts",
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

			items, nextToken, err := c.List(context.Background(), "/storage_accounts", params)
			if err != nil {
				return fmt.Errorf("could not list storage accounts: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, storageAccountListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newStorageAccountsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a storage account by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/storage_accounts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get storage account: %w", err)
			}

			return output.PrintItem(cmd, result, storageAccountDetailColumns)
		},
	}
}

func newStorageAccountsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new storage account",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			customerID, _ := cmd.Flags().GetString("customer-id")
			body := map[string]interface{}{
				"customer_id": customerID,
			}

			result, err := c.Post(context.Background(), "/storage_accounts", body)
			if err != nil {
				return fmt.Errorf("could not create storage account: %w", err)
			}

			return output.PrintItem(cmd, result, storageAccountDetailColumns)
		},
	}
	cmd.Flags().String("customer-id", "", "Customer ID")
	_ = cmd.MarkFlagRequired("customer-id")
	return cmd
}

func newStorageAccountsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a storage account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/storage_accounts/"+args[0]); err != nil {
				return fmt.Errorf("could not delete storage account: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Storage account %s deleted.\n", args[0])
			return nil
		},
	}
}
