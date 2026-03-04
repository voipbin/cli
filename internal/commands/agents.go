package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetAgentsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetAgentsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list agents: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, agentListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetAgentsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get agent: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
		},
	}
}

func newAgentsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			detail, _ := cmd.Flags().GetString("detail")
			ringMethod, _ := cmd.Flags().GetString("ring-method")
			permission, _ := cmd.Flags().GetInt("permission")

			body := voipbin_client.PostAgentsJSONRequestBody{
				Name:       name,
				Username:   username,
				Password:   password,
				Detail:     detail,
				RingMethod: voipbin_client.AgentManagerAgentRingMethod(ringMethod),
				Permission: voipbin_client.AgentManagerAgentPermission(permission),
				Addresses:  []voipbin_client.CommonAddress{},
				TagIds:     []string{},
			}

			resp, err := client.PostAgentsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create agent: %w", err)
			}
			if resp.StatusCode() != 201 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON201 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON201, agentDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			ringMethod, _ := cmd.Flags().GetString("ring-method")

			body := voipbin_client.PutAgentsIdJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}
			if ringMethod != "" {
				rm := voipbin_client.AgentManagerAgentRingMethod(ringMethod)
				body.RingMethod = &rm
			}

			resp, err := client.PutAgentsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update agent: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteAgentsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete agent: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			body := voipbin_client.PutAgentsIdAddressesJSONRequestBody{
				Addresses: &[]voipbin_client.CommonAddress{},
			}

			resp, err := client.PutAgentsIdAddressesWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update addresses: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
		},
	}
	return cmd
}

func newAgentsUpdatePasswordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-password <id>",
		Short: "Update agent password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			password, _ := cmd.Flags().GetString("password")
			body := voipbin_client.PutAgentsIdPasswordJSONRequestBody{
				Password: &password,
			}

			resp, err := client.PutAgentsIdPasswordWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update password: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			permission, _ := cmd.Flags().GetInt("permission")
			perm := voipbin_client.AgentManagerAgentPermission(permission)
			body := voipbin_client.PutAgentsIdPermissionJSONRequestBody{
				Permission: &perm,
			}

			resp, err := client.PutAgentsIdPermissionWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update permission: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")
			s := voipbin_client.AgentManagerAgentStatus(status)
			body := voipbin_client.PutAgentsIdStatusJSONRequestBody{
				Status: &s,
			}

			resp, err := client.PutAgentsIdStatusWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update status: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			tagIDs, _ := cmd.Flags().GetStringArray("tag-id")
			body := voipbin_client.PutAgentsIdTagIdsJSONRequestBody{
				TagIds: &tagIDs,
			}

			resp, err := client.PutAgentsIdTagIdsWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update tag IDs: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, agentDetailColumns)
		},
	}
	cmd.Flags().StringArray("tag-id", []string{}, "Tag ID (repeatable)")
	return cmd
}
