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

func newRagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rags",
		Short: "Manage RAG knowledge bases",
	}
	cmd.AddCommand(
		newRagsListCmd(),
		newRagsGetCmd(),
		newRagsCreateCmd(),
		newRagsUpdateCmd(),
		newRagsDeleteCmd(),
		newRagsAddSourceCmd(),
	)
	return cmd
}

var ragListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DESCRIPTION", Field: "description"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var ragDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DESCRIPTION", Field: "description"},
	{Name: "STATUS", Field: "status"},
	{Name: "SOURCES", Field: "sources"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newRagsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List RAG knowledge bases",
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

			items, nextToken, err := c.List(context.Background(), "/rags", params)
			if err != nil {
				return fmt.Errorf("could not list rags: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, ragListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newRagsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a RAG knowledge base by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/rags/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get rag: %w", err)
			}

			return output.PrintItem(cmd, result, ragDetailColumns)
		},
	}
}

func newRagsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new RAG knowledge base",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")

			body := map[string]interface{}{
				"name":        name,
				"description": description,
			}

			result, err := c.Post(context.Background(), "/rags", body)
			if err != nil {
				return fmt.Errorf("could not create rag: %w", err)
			}

			return output.PrintItem(cmd, result, ragDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "RAG name")
	cmd.Flags().String("description", "", "RAG description")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newRagsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a RAG knowledge base",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if description != "" {
				body["description"] = description
			}

			result, err := c.Put(context.Background(), "/rags/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update rag: %w", err)
			}

			return output.PrintItem(cmd, result, ragDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "RAG name")
	cmd.Flags().String("description", "", "RAG description")
	return cmd
}

func newRagsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a RAG knowledge base",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/rags/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete rag: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "RAG %s deleted.\n", args[0])
			return nil
		},
	}
}

func newRagsAddSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-source <id>",
		Short: "Add sources to a RAG knowledge base",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			storageFileIDs, _ := cmd.Flags().GetStringSlice("storage-file-ids")
			sourceURLs, _ := cmd.Flags().GetStringSlice("source-urls")

			if len(storageFileIDs) == 0 && len(sourceURLs) == 0 {
				return fmt.Errorf("at least one of --storage-file-ids or --source-urls is required")
			}

			body := map[string]interface{}{}
			if len(storageFileIDs) > 0 {
				body["storage_file_ids"] = storageFileIDs
			}
			if len(sourceURLs) > 0 {
				body["source_urls"] = sourceURLs
			}

			result, err := c.Post(context.Background(), "/rags/"+args[0]+"/sources", body)
			if err != nil {
				return fmt.Errorf("could not add sources: %w", err)
			}

			return output.PrintItem(cmd, result, ragDetailColumns)
		},
	}
	cmd.Flags().StringSlice("storage-file-ids", nil, "Storage file IDs to add as sources")
	cmd.Flags().StringSlice("source-urls", nil, "URLs to add as sources")
	return cmd
}
