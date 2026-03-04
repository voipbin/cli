package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetTagsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetTagsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list tags: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, tagListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetTagsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get tag: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, tagDetailColumns)
		},
	}
}

func newTagsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new tag",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PostTagsJSONRequestBody{
				Name:   name,
				Detail: detail,
			}

			resp, err := client.PostTagsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create tag: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, tagDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PutTagsIdJSONRequestBody{
				Name:   name,
				Detail: detail,
			}

			resp, err := client.PutTagsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update tag: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, tagDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteTagsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete tag: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Tag %s deleted.\n", args[0])
			return nil
		},
	}
}
