package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
)

func newAgentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Manage agents",
	}
	cmd.AddCommand(
		newAgentsListCmd(),
		newAgentsGetCmd(),
		newAgentsCreateCmd(),
		newAgentsUpdateCmd(),
		newAgentsDeleteCmd(),
		newAgentsUpdateAddressesCmd(),
		newAgentsUpdatePasswordCmd(),
		newAgentsUpdatePermissionCmd(),
		newAgentsUpdateStatusCmd(),
		newAgentsUpdateTagIdsCmd(),
	)
	return cmd
}

var agentListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "USERNAME", Field: "username"},
	{Name: "STATUS", Field: "status"},
	{Name: "PERMISSION", Field: "permission"},
	{Name: "CREATED", Field: "tm_create"},
}

var agentDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "USERNAME", Field: "username"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "STATUS", Field: "status"},
	{Name: "PERMISSION", Field: "permission"},
	{Name: "RING_METHOD", Field: "ring_method"},
	{Name: "CREATED", Field: "tm_create"},
}

func newAgentsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List agents",
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

			items, nextToken, err := c.List(context.Background(), "/agents", params)
			if err != nil {
				return fmt.Errorf("could not list agents: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, agentListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newAgentsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an agent by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/agents/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get agent: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
}

func newAgentsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			detail, _ := cmd.Flags().GetString("detail")
			ringMethod, _ := cmd.Flags().GetString("ring-method")
			permission, _ := cmd.Flags().GetInt("permission")

			body := map[string]interface{}{
				"name":        name,
				"username":    username,
				"password":    password,
				"detail":      detail,
				"ring_method": ringMethod,
				"permission":  permission,
				"addresses":   []interface{}{},
				"tag_ids":     []interface{}{},
			}

			result, err := c.Post(context.Background(), "/agents", body)
			if err != nil {
				return fmt.Errorf("could not create agent: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Agent name")
	cmd.Flags().String("username", "", "Agent username")
	cmd.Flags().String("password", "", "Agent password")
	cmd.Flags().String("detail", "", "Agent detail")
	cmd.Flags().String("ring-method", "ringall", "Ring method")
	cmd.Flags().Int("permission", 0, "Permission level")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func newAgentsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			ringMethod, _ := cmd.Flags().GetString("ring-method")

			// Pattern 2 (conditional): only add to map if value is non-empty.
			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if ringMethod != "" {
				body["ring_method"] = ringMethod
			}

			result, err := c.Put(context.Background(), "/agents/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update agent: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Agent name")
	cmd.Flags().String("detail", "", "Agent detail")
	cmd.Flags().String("ring-method", "", "Ring method")
	return cmd
}

func newAgentsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an agent",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/agents/"+args[0]); err != nil {
				return fmt.Errorf("could not delete agent: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Agent %s deleted.\n", args[0])
			return nil
		},
	}
}

func newAgentsUpdateAddressesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-addresses <id>",
		Short: "Update agent addresses",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			addressesJSON, _ := cmd.Flags().GetString("addresses")
			var addresses interface{}
			if err := json.Unmarshal([]byte(addressesJSON), &addresses); err != nil {
				return fmt.Errorf("invalid addresses JSON: %w", err)
			}

			body := map[string]interface{}{
				"addresses": addresses,
			}

			result, err := c.Put(context.Background(), "/agents/"+args[0]+"/addresses", body)
			if err != nil {
				return fmt.Errorf("could not update addresses: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().String("addresses", "", "Addresses as JSON array (e.g. '[{\"type\":\"sip\",\"target\":\"user@host\"}]')")
	_ = cmd.MarkFlagRequired("addresses")
	return cmd
}

func newAgentsUpdatePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-password <id>",
		Short: "Update agent password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			password, _ := cmd.Flags().GetString("password")
			body := map[string]interface{}{
				"password": password,
			}

			result, err := c.Put(context.Background(), "/agents/"+args[0]+"/password", body)
			if err != nil {
				return fmt.Errorf("could not update password: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().String("password", "", "New password")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func newAgentsUpdatePermissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-permission <id>",
		Short: "Update agent permission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			permission, _ := cmd.Flags().GetInt("permission")
			body := map[string]interface{}{
				"permission": permission,
			}

			result, err := c.Put(context.Background(), "/agents/"+args[0]+"/permission", body)
			if err != nil {
				return fmt.Errorf("could not update permission: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().Int("permission", 0, "Permission level")
	_ = cmd.MarkFlagRequired("permission")
	return cmd
}

func newAgentsUpdateStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-status <id>",
		Short: "Update agent status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")
			body := map[string]interface{}{
				"status": status,
			}

			result, err := c.Put(context.Background(), "/agents/"+args[0]+"/status", body)
			if err != nil {
				return fmt.Errorf("could not update status: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().String("status", "", "Agent status")
	_ = cmd.MarkFlagRequired("status")
	return cmd
}

func newAgentsUpdateTagIdsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-tag-ids <id>",
		Short: "Update agent tag IDs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			tagIDs, _ := cmd.Flags().GetStringArray("tag-id")
			body := map[string]interface{}{
				"tag_ids": tagIDs,
			}

			result, err := c.Put(context.Background(), "/agents/"+args[0]+"/tag_ids", body)
			if err != nil {
				return fmt.Errorf("could not update tag IDs: %w", err)
			}

			return output.PrintItem(cmd, result, agentDetailColumns)
		},
	}
	cmd.Flags().StringArray("tag-id", []string{}, "Tag ID (repeatable)")
	_ = cmd.MarkFlagRequired("tag-id")
	return cmd
}
