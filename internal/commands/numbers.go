package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetNumbersParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetNumbersWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list numbers: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, numberListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetNumbersIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get number: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, numberDetailColumns)
		},
	}
}

func newNumbersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create (purchase) a phone number",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			callFlowID, _ := cmd.Flags().GetString("call-flow-id")
			messageFlowID, _ := cmd.Flags().GetString("message-flow-id")

			body := voipbin_client.PostNumbersJSONRequestBody{
				Name:          name,
				Detail:        detail,
				CallFlowId:    callFlowID,
				MessageFlowId: messageFlowID,
			}

			resp, err := client.PostNumbersWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create number: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, numberDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			callFlowID, _ := cmd.Flags().GetString("call-flow-id")
			messageFlowID, _ := cmd.Flags().GetString("message-flow-id")

			body := voipbin_client.PutNumbersIdJSONRequestBody{
				Name:          name,
				Detail:        detail,
				CallFlowId:    callFlowID,
				MessageFlowId: messageFlowID,
			}

			resp, err := client.PutNumbersIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update number: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, numberDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteNumbersIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete number: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			tmRenew, _ := cmd.Flags().GetString("tm-renew")
			body := voipbin_client.PostNumbersRenewJSONRequestBody{
				TmRenew: tmRenew,
			}

			resp, err := client.PostNumbersRenewWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not renew number: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			callFlowID, _ := cmd.Flags().GetString("call-flow-id")
			messageFlowID, _ := cmd.Flags().GetString("message-flow-id")

			body := voipbin_client.PutNumbersIdFlowIdsJSONRequestBody{
				CallFlowId:    callFlowID,
				MessageFlowId: messageFlowID,
			}

			resp, err := client.PutNumbersIdFlowIdsWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update flow IDs: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, numberDetailColumns)
		},
	}
	cmd.Flags().String("call-flow-id", "", "Call flow ID")
	cmd.Flags().String("message-flow-id", "", "Message flow ID")
	return cmd
}
