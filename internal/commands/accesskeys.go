package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newAccesskeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accesskeys",
		Short: "Manage access keys",
	}
	cmd.AddCommand(
		newAccesskeysListCmd(),
		newAccesskeysGetCmd(),
		newAccesskeysCreateCmd(),
		newAccesskeysUpdateCmd(),
		newAccesskeysDeleteCmd(),
	)
	return cmd
}

var accesskeyListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "TM_EXPIRE", Field: "tm_expire"},
	{Name: "CREATED", Field: "tm_create"},
}

var accesskeyDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "TOKEN", Field: "token"},
	{Name: "TM_EXPIRE", Field: "tm_expire"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newAccesskeysListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List access keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetAccesskeysParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetAccesskeysWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list access keys: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, accesskeyListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newAccesskeysGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an access key by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetAccesskeysIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get access key: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, accesskeyDetailColumns)
		},
	}
}

func newAccesskeysCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new access key",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			expire, _ := cmd.Flags().GetInt("expire")

			body := voipbin_client.PostAccesskeysJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}
			if expire > 0 {
				body.Expire = &expire
			}

			resp, err := client.PostAccesskeysWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create access key: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, accesskeyDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Access key name")
	cmd.Flags().String("detail", "", "Access key detail")
	cmd.Flags().Int("expire", 0, "Expiry (seconds from now)")
	return cmd
}

func newAccesskeysUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an access key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PutAccesskeysIdJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}

			resp, err := client.PutAccesskeysIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update access key: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, accesskeyDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Access key name")
	cmd.Flags().String("detail", "", "Access key detail")
	return cmd
}

func newAccesskeysDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an access key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteAccesskeysIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete access key: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Access key %s deleted.\n", args[0])
			return nil
		},
	}
}
