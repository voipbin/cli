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

			items, nextToken, err := c.List(context.Background(), "/conversation_accounts", params)
			if err != nil {
				return fmt.Errorf("could not list conversation accounts: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, conversationAccountListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/conversation_accounts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get conversation account: %w", err)
			}

			return output.PrintItem(cmd, result, conversationAccountDetailColumns)
		},
	}
}

func newConversationAccountsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new conversation account",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			accountType, _ := cmd.Flags().GetString("type")
			token, _ := cmd.Flags().GetString("token")
			secret, _ := cmd.Flags().GetString("secret")

			body := map[string]interface{}{
				"name":   name,
				"detail": detail,
				"type":   accountType,
				"token":  token,
				"secret": secret,
			}

			result, err := c.Post(context.Background(), "/conversation_accounts", body)
			if err != nil {
				return fmt.Errorf("could not create conversation account: %w", err)
			}

			return output.PrintItem(cmd, result, conversationAccountDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			token, _ := cmd.Flags().GetString("token")
			secret, _ := cmd.Flags().GetString("secret")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if token != "" {
				body["token"] = token
			}
			if secret != "" {
				body["secret"] = secret
			}

			result, err := c.Put(context.Background(), "/conversation_accounts/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update conversation account: %w", err)
			}

			return output.PrintItem(cmd, result, conversationAccountDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/conversation_accounts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete conversation account: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conversation account %s deleted.\n", args[0])
			return nil
		},
	}
}
