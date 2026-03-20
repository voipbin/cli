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

func newFlowsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flows",
		Short: "Manage flows",
	}
	cmd.AddCommand(
		newFlowsListCmd(),
		newFlowsGetCmd(),
		newFlowsCreateCmd(),
		newFlowsUpdateCmd(),
		newFlowsDeleteCmd(),
	)
	return cmd
}

var flowListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "TYPE", Field: "type"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CREATED", Field: "tm_create"},
}

var flowDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "TYPE", Field: "type"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "ON_COMPLETE_FLOW_ID", Field: "on_complete_flow_id"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newFlowsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List flows",
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
			items, nextToken, err := c.List(context.Background(), "/flows", params)
			if err != nil {
				return fmt.Errorf("could not list flows: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, flowListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newFlowsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a flow by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/flows/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get flow: %w", err)
			}
			return output.PrintItem(cmd, item, flowDetailColumns)
		},
	}
}

func newFlowsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new flow",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			onCompleteFlowID, _ := cmd.Flags().GetString("on-complete-flow-id")
			body := map[string]interface{}{
				"name":    name,
				"detail":  detail,
				"actions": []interface{}{},
			}
			if onCompleteFlowID != "" {
				body["on_complete_flow_id"] = onCompleteFlowID
			}
			item, err := c.Post(context.Background(), "/flows", body)
			if err != nil {
				return fmt.Errorf("could not create flow: %w", err)
			}
			return output.PrintItem(cmd, item, flowDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Flow name")
	cmd.Flags().String("detail", "", "Flow detail")
	cmd.Flags().String("on-complete-flow-id", "", "Flow ID to execute on completion")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newFlowsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			onCompleteFlowID, _ := cmd.Flags().GetString("on-complete-flow-id")
			body := map[string]interface{}{
				"name":    name,
				"detail":  detail,
				"actions": []interface{}{},
			}
			if onCompleteFlowID != "" {
				body["on_complete_flow_id"] = onCompleteFlowID
			}
			item, err := c.Put(context.Background(), "/flows/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update flow: %w", err)
			}
			return output.PrintItem(cmd, item, flowDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("on-complete-flow-id", "", "Flow ID to execute on completion")
	return cmd
}

func newFlowsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a flow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/flows/"+args[0]); err != nil {
				return fmt.Errorf("could not delete flow: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Flow %s deleted.\n", args[0])
			return nil
		},
	}
}
