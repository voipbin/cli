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

func newNumbersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "numbers",
		Short: "Manage phone numbers",
	}
	cmd.AddCommand(
		newNumbersListCmd(),
		newNumbersGetCmd(),
		newNumbersCreateCmd(),
		newNumbersUpdateCmd(),
		newNumbersDeleteCmd(),
		newNumbersRenewCmd(),
		newNumbersUpdateFlowIDsCmd(),
	)
	return cmd
}

var numberListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NUMBER", Field: "number"},
	{Name: "NAME", Field: "name"},
	{Name: "STATUS", Field: "status"},
	{Name: "CALL_FLOW_ID", Field: "call_flow_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var numberDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NUMBER", Field: "number"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "STATUS", Field: "status"},
	{Name: "CALL_FLOW_ID", Field: "call_flow_id"},
	{Name: "MESSAGE_FLOW_ID", Field: "message_flow_id"},
	{Name: "T38_ENABLED", Field: "t38_enabled"},
	{Name: "EMERGENCY_ENABLED", Field: "emergency_enabled"},
	{Name: "CREATED", Field: "tm_create"},
}

func newNumbersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List phone numbers",
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

			items, nextToken, err := c.List(context.Background(), "/numbers", params)
			if err != nil {
				return fmt.Errorf("could not list numbers: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, numberListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newNumbersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a phone number by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/numbers/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get number: %w", err)
			}

			return output.PrintItem(cmd, result, numberDetailColumns)
		},
	}
}

func newNumbersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create (purchase) a phone number",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			callFlowID, _ := cmd.Flags().GetString("call-flow-id")
			messageFlowID, _ := cmd.Flags().GetString("message-flow-id")

			body := map[string]interface{}{
				"name":            name,
				"detail":          detail,
				"call_flow_id":    callFlowID,
				"message_flow_id": messageFlowID,
			}

			result, err := c.Post(context.Background(), "/numbers", body)
			if err != nil {
				return fmt.Errorf("could not create number: %w", err)
			}

			return output.PrintItem(cmd, result, numberDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Number name")
	cmd.Flags().String("detail", "", "Number detail")
	cmd.Flags().String("call-flow-id", "", "Call flow ID")
	cmd.Flags().String("message-flow-id", "", "Message flow ID")
	return cmd
}

func newNumbersUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a phone number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			callFlowID, _ := cmd.Flags().GetString("call-flow-id")
			messageFlowID, _ := cmd.Flags().GetString("message-flow-id")

			body := map[string]interface{}{
				"name":            name,
				"detail":          detail,
				"call_flow_id":    callFlowID,
				"message_flow_id": messageFlowID,
			}

			result, err := c.Put(context.Background(), "/numbers/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update number: %w", err)
			}

			return output.PrintItem(cmd, result, numberDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Number name")
	cmd.Flags().String("detail", "", "Number detail")
	cmd.Flags().String("call-flow-id", "", "Call flow ID")
	cmd.Flags().String("message-flow-id", "", "Message flow ID")
	return cmd
}

func newNumbersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a phone number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/numbers/"+args[0]); err != nil {
				return fmt.Errorf("could not delete number: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Number %s deleted.\n", args[0])
			return nil
		},
	}
}

func newNumbersRenewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew",
		Short: "Renew a phone number",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			tmRenew, _ := cmd.Flags().GetString("tm-renew")
			body := map[string]interface{}{
				"tm_renew": tmRenew,
			}

			if _, err := c.Post(context.Background(), "/numbers/renew", body); err != nil {
				return fmt.Errorf("could not renew number: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Number renewal scheduled.\n")
			return nil
		},
	}
	cmd.Flags().String("tm-renew", "", "Renewal timestamp")
	_ = cmd.MarkFlagRequired("tm-renew")
	return cmd
}

func newNumbersUpdateFlowIDsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-flow-ids <id>",
		Short: "Update flow IDs for a phone number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			callFlowID, _ := cmd.Flags().GetString("call-flow-id")
			messageFlowID, _ := cmd.Flags().GetString("message-flow-id")

			body := map[string]interface{}{
				"call_flow_id":    callFlowID,
				"message_flow_id": messageFlowID,
			}

			result, err := c.Put(context.Background(), "/numbers/"+args[0]+"/flow_ids", body)
			if err != nil {
				return fmt.Errorf("could not update flow IDs: %w", err)
			}

			return output.PrintItem(cmd, result, numberDetailColumns)
		},
	}
	cmd.Flags().String("call-flow-id", "", "Call flow ID")
	cmd.Flags().String("message-flow-id", "", "Message flow ID")
	return cmd
}
