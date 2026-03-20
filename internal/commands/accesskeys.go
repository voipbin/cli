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

			items, nextToken, err := c.List(context.Background(), "/accesskeys", params)
			if err != nil {
				return fmt.Errorf("could not list access keys: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, accesskeyListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/accesskeys/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get access key: %w", err)
			}

			return output.PrintItem(cmd, result, accesskeyDetailColumns)
		},
	}
}

func newAccesskeysCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new access key",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			expire, _ := cmd.Flags().GetInt("expire")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if expire > 0 {
				body["expire"] = expire
			}

			result, err := c.Post(context.Background(), "/accesskeys", body)
			if err != nil {
				return fmt.Errorf("could not create access key: %w", err)
			}

			return output.PrintItem(cmd, result, accesskeyDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}

			result, err := c.Put(context.Background(), "/accesskeys/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update access key: %w", err)
			}

			return output.PrintItem(cmd, result, accesskeyDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/accesskeys/"+args[0]); err != nil {
				return fmt.Errorf("could not delete access key: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Access key %s deleted.\n", args[0])
			return nil
		},
	}
}
