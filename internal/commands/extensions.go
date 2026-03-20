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

func newExtensionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extensions",
		Short: "Manage extensions",
	}
	cmd.AddCommand(
		newExtensionsListCmd(),
		newExtensionsGetCmd(),
		newExtensionsCreateCmd(),
		newExtensionsUpdateCmd(),
		newExtensionsDeleteCmd(),
	)
	return cmd
}

var extensionListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "EXTENSION", Field: "extension"},
	{Name: "USERNAME", Field: "username"},
	{Name: "CREATED", Field: "tm_create"},
}

var extensionDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "EXTENSION", Field: "extension"},
	{Name: "USERNAME", Field: "username"},
	{Name: "DOMAIN_NAME", Field: "domain_name"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newExtensionsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List extensions",
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
			items, nextToken, err := c.List(context.Background(), "/extensions", params)
			if err != nil {
				return fmt.Errorf("could not list extensions: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, extensionListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newExtensionsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an extension by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/extensions/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get extension: %w", err)
			}
			return output.PrintItem(cmd, item, extensionDetailColumns)
		},
	}
}

func newExtensionsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			extension, _ := cmd.Flags().GetString("extension")
			password, _ := cmd.Flags().GetString("password")
			body := map[string]interface{}{
				"name":      name,
				"detail":    detail,
				"extension": extension,
				"password":  password,
			}
			item, err := c.Post(context.Background(), "/extensions", body)
			if err != nil {
				return fmt.Errorf("could not create extension: %w", err)
			}
			return output.PrintItem(cmd, item, extensionDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Extension name")
	cmd.Flags().String("detail", "", "Extension detail")
	cmd.Flags().String("extension", "", "Extension number")
	cmd.Flags().String("password", "", "Extension password")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("extension")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func newExtensionsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an extension",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			password, _ := cmd.Flags().GetString("password")
			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if password != "" {
				body["password"] = password
			}
			item, err := c.Put(context.Background(), "/extensions/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update extension: %w", err)
			}
			return output.PrintItem(cmd, item, extensionDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("password", "", "New password")
	return cmd
}

func newExtensionsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an extension",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/extensions/"+args[0]); err != nil {
				return fmt.Errorf("could not delete extension: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Extension %s deleted.\n", args[0])
			return nil
		},
	}
}
