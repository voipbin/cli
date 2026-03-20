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

func newConversationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversations",
		Short: "Manage conversations",
	}
	cmd.AddCommand(
		newConversationsListCmd(),
		newConversationsGetCmd(),
		newConversationsUpdateCmd(),
		newConversationsListMessagesCmd(),
		newConversationsCreateMessageCmd(),
	)
	return cmd
}

var conversationListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "ACCOUNT_ID", Field: "account_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "CREATED", Field: "tm_create"},
}

var conversationDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "ACCOUNT_ID", Field: "account_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "OWNER_TYPE", Field: "owner_type"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

var conversationMessageColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CONVERSATION_ID", Field: "conversation_id"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "STATUS", Field: "status"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
}

func newConversationsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List conversations",
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
			items, nextToken, err := c.List(context.Background(), "/conversations", params)
			if err != nil {
				return fmt.Errorf("could not list conversations: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, conversationListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newConversationsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a conversation by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/conversations/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get conversation: %w", err)
			}
			return output.PrintItem(cmd, item, conversationDetailColumns)
		},
	}
}

func newConversationsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			ownerID, _ := cmd.Flags().GetString("owner-id")
			ownerType, _ := cmd.Flags().GetString("owner-type")
			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if ownerID != "" {
				body["owner_id"] = ownerID
			}
			if ownerType != "" {
				body["owner_type"] = ownerType
			}
			item, err := c.Put(context.Background(), "/conversations/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update conversation: %w", err)
			}
			return output.PrintItem(cmd, item, conversationDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("owner-id", "", "New owner ID")
	cmd.Flags().String("owner-type", "", "New owner type")
	return cmd
}

func newConversationsListMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-messages <id>",
		Short: "List messages in a conversation",
		Args:  cobra.ExactArgs(1),
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
			items, nextToken, err := c.List(context.Background(), "/conversations/"+args[0]+"/messages", params)
			if err != nil {
				return fmt.Errorf("could not list messages: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, conversationMessageColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newConversationsCreateMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-message <id>",
		Short: "Create a message in a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			text, _ := cmd.Flags().GetString("text")
			body := map[string]interface{}{
				"text":   text,
				"medias": []interface{}{},
			}
			item, err := c.Post(context.Background(), "/conversations/"+args[0]+"/messages", body)
			if err != nil {
				return fmt.Errorf("could not create message: %w", err)
			}
			return output.PrintItem(cmd, item, conversationMessageColumns)
		},
	}
	cmd.Flags().String("text", "", "Message text")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}
