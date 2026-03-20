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

func newTrunksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trunks",
		Short: "Manage trunks",
	}
	cmd.AddCommand(
		newTrunksListCmd(),
		newTrunksGetCmd(),
		newTrunksCreateCmd(),
		newTrunksUpdateCmd(),
		newTrunksDeleteCmd(),
	)
	return cmd
}

var trunkListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DOMAIN_NAME", Field: "domain_name"},
	{Name: "USERNAME", Field: "username"},
	{Name: "CREATED", Field: "tm_create"},
}

var trunkDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "DOMAIN_NAME", Field: "domain_name"},
	{Name: "USERNAME", Field: "username"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newTrunksListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List trunks",
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

			items, nextToken, err := c.List(context.Background(), "/trunks", params)
			if err != nil {
				return fmt.Errorf("could not list trunks: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, trunkListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newTrunksGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a trunk by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/trunks/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get trunk: %w", err)
			}

			return output.PrintItem(cmd, result, trunkDetailColumns)
		},
	}
}

func newTrunksCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new trunk",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			domainName, _ := cmd.Flags().GetString("domain-name")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			detail, _ := cmd.Flags().GetString("detail")

			body := map[string]interface{}{
				"name":        name,
				"domain_name": domainName,
				"username":    username,
				"password":    password,
				"detail":      detail,
				"allowed_ips": []interface{}{},
				"auth_types":  []interface{}{},
			}

			result, err := c.Post(context.Background(), "/trunks", body)
			if err != nil {
				return fmt.Errorf("could not create trunk: %w", err)
			}

			return output.PrintItem(cmd, result, trunkDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Trunk name")
	cmd.Flags().String("domain-name", "", "Domain name")
	cmd.Flags().String("username", "", "Username")
	cmd.Flags().String("password", "", "Password")
	cmd.Flags().String("detail", "", "Trunk detail")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newTrunksUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a trunk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			detail, _ := cmd.Flags().GetString("detail")

			body := map[string]interface{}{
				"name":        name,
				"username":    username,
				"password":    password,
				"detail":      detail,
				"allowed_ips": []interface{}{},
				"auth_types":  []interface{}{},
			}

			result, err := c.Put(context.Background(), "/trunks/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update trunk: %w", err)
			}

			return output.PrintItem(cmd, result, trunkDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Trunk name")
	cmd.Flags().String("username", "", "Username")
	cmd.Flags().String("password", "", "Password")
	cmd.Flags().String("detail", "", "Trunk detail")
	return cmd
}

func newTrunksDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a trunk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/trunks/"+args[0]); err != nil {
				return fmt.Errorf("could not delete trunk: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Trunk %s deleted.\n", args[0])
			return nil
		},
	}
}
