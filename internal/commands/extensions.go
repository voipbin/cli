package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetExtensionsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetExtensionsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list extensions: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, extensionListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetExtensionsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get extension: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, extensionDetailColumns)
		},
	}
}

func newExtensionsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			extension, _ := cmd.Flags().GetString("extension")
			password, _ := cmd.Flags().GetString("password")
			body := voipbin_client.PostExtensionsJSONRequestBody{
				Name:      name,
				Detail:    detail,
				Extension: extension,
				Password:  password,
			}
			resp, err := client.PostExtensionsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create extension: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, extensionDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			password, _ := cmd.Flags().GetString("password")
			body := voipbin_client.PutExtensionsIdJSONRequestBody{
				Name:     name,
				Detail:   detail,
				Password: password,
			}
			resp, err := client.PutExtensionsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update extension: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, extensionDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteExtensionsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete extension: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Extension %s deleted.\n", args[0])
			return nil
		},
	}
}
