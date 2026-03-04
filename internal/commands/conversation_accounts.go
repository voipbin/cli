package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newConversationAccountsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversation-accounts",
		Short: "Manage conversation accounts",
	}
	cmd.AddCommand(
		newConversationAccountsListCmd(),
		newConversationAccountsGetCmd(),
		newConversationAccountsCreateCmd(),
		newConversationAccountsUpdateCmd(),
		newConversationAccountsDeleteCmd(),
	)
	return cmd
}

var conversationAccountListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "TYPE", Field: "type"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CREATED", Field: "tm_create"},
}

var conversationAccountDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "TYPE", Field: "type"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newConversationAccountsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List conversation accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetConversationAccountsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetConversationAccountsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list conversation accounts: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, conversationAccountListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newConversationAccountsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a conversation account by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetConversationAccountsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get conversation account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, conversationAccountDetailColumns)
		},
	}
}

func newConversationAccountsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new conversation account",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			accountType, _ := cmd.Flags().GetString("type")
			token, _ := cmd.Flags().GetString("token")
			secret, _ := cmd.Flags().GetString("secret")
			body := voipbin_client.PostConversationAccountsJSONRequestBody{
				Name:   name,
				Detail: detail,
				Type:   voipbin_client.ConversationManagerAccountType(accountType),
				Token:  token,
				Secret: secret,
			}
			resp, err := client.PostConversationAccountsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create conversation account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, conversationAccountDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Account name")
	cmd.Flags().String("detail", "", "Account detail")
	cmd.Flags().String("type", "", "Account type")
	cmd.Flags().String("token", "", "API token")
	cmd.Flags().String("secret", "", "API secret")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newConversationAccountsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a conversation account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			token, _ := cmd.Flags().GetString("token")
			secret, _ := cmd.Flags().GetString("secret")
			body := voipbin_client.PutConversationAccountsIdJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}
			if token != "" {
				body.Token = &token
			}
			if secret != "" {
				body.Secret = &secret
			}
			resp, err := client.PutConversationAccountsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update conversation account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, conversationAccountDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("token", "", "New token")
	cmd.Flags().String("secret", "", "New secret")
	return cmd
}

func newConversationAccountsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a conversation account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteConversationAccountsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete conversation account: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Conversation account %s deleted.\n", args[0])
			return nil
		},
	}
}
