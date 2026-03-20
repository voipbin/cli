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

func newTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Manage tags",
	}
	cmd.AddCommand(
		newTagsListCmd(),
		newTagsGetCmd(),
		newTagsCreateCmd(),
		newTagsUpdateCmd(),
		newTagsDeleteCmd(),
	)
	return cmd
}

var tagListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CREATED", Field: "tm_create"},
}

var tagDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newTagsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tags",
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

			items, nextToken, err := c.List(context.Background(), "/tags", params)
			if err != nil {
				return fmt.Errorf("could not list tags: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, tagListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newTagsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a tag by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/tags/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get tag: %w", err)
			}

			return output.PrintItem(cmd, result, tagDetailColumns)
		},
	}
}

func newTagsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new tag",
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

			result, err := c.Post(context.Background(), "/tags", body)
			if err != nil {
				return fmt.Errorf("could not create tag: %w", err)
			}

			return output.PrintItem(cmd, result, tagDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Tag name")
	cmd.Flags().String("detail", "", "Tag detail")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newTagsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a tag",
		Args:  cobra.ExactArgs(1),
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

			result, err := c.Put(context.Background(), "/tags/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update tag: %w", err)
			}

			return output.PrintItem(cmd, result, tagDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Tag name")
	cmd.Flags().String("detail", "", "Tag detail")
	return cmd
}

func newTagsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/tags/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete tag: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Tag %s deleted.\n", args[0])
			return nil
		},
	}
}
