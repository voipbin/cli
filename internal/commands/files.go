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

func newFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files",
		Short: "Manage files",
	}
	cmd.AddCommand(
		newFilesListCmd(),
		newFilesGetCmd(),
		newFilesDeleteCmd(),
	)
	return cmd
}

var fileListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "FILENAME", Field: "filename"},
	{Name: "FILESIZE", Field: "filesize"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var fileDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "FILENAME", Field: "filename"},
	{Name: "FILESIZE", Field: "filesize"},
	{Name: "CREATED", Field: "tm_create"},
}

func newFilesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files",
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

			items, nextToken, err := c.List(context.Background(), "/files", params)
			if err != nil {
				return fmt.Errorf("could not list files: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, fileListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newFilesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a file by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/files/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get file: %w", err)
			}

			return output.PrintItem(cmd, result, fileDetailColumns)
		},
	}
}

func newFilesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/files/"+args[0]); err != nil {
				return fmt.Errorf("could not delete file: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "File %s deleted.\n", args[0])
			return nil
		},
	}
}
