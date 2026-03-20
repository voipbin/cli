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

func newStorageFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "storage-files",
		Short: "Manage storage files",
	}
	cmd.AddCommand(
		newStorageFilesListCmd(),
		newStorageFilesGetCmd(),
		newStorageFilesDeleteCmd(),
	)
	return cmd
}

var storageFileListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "FILENAME", Field: "filename"},
	{Name: "FILESIZE", Field: "filesize"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var storageFileDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "FILENAME", Field: "filename"},
	{Name: "FILESIZE", Field: "filesize"},
	{Name: "CREATED", Field: "tm_create"},
}

func newStorageFilesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List storage files",
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

			items, nextToken, err := c.List(context.Background(), "/storage_files", params)
			if err != nil {
				return fmt.Errorf("could not list storage files: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, storageFileListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newStorageFilesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a storage file by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/storage_files/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get storage file: %w", err)
			}

			return output.PrintItem(cmd, result, storageFileDetailColumns)
		},
	}
}

func newStorageFilesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a storage file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/storage_files/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete storage file: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Storage file %s deleted.\n", args[0])
			return nil
		},
	}
}
