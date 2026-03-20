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

func newChatmessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chatmessages",
		Short: "Manage chat messages",
	}
	cmd.AddCommand(
		newChatmessagesListCmd(),
		newChatmessagesGetCmd(),
		newChatmessagesCreateCmd(),
		newChatmessagesDeleteCmd(),
	)
	return cmd
}

var chatmessageListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CHAT_ID", Field: "chat_id"},
	{Name: "TYPE", Field: "type"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
}

var chatmessageDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CHAT_ID", Field: "chat_id"},
	{Name: "TYPE", Field: "type"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newChatmessagesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List chat messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			chatID, _ := cmd.Flags().GetString("chat-id")
			params := url.Values{}
			params.Set("chat_id", chatID)
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}
			items, nextToken, err := c.List(context.Background(), "/chatmessages", params)
			if err != nil {
				return fmt.Errorf("could not list chatmessages: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, chatmessageListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	cmd.Flags().String("chat-id", "", "Chat ID to list messages for")
	_ = cmd.MarkFlagRequired("chat-id")
	return cmd
}

func newChatmessagesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a chat message by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/chatmessages/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get chatmessage: %w", err)
			}
			return output.PrintItem(cmd, item, chatmessageDetailColumns)
		},
	}
}

func newChatmessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chat message",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			chatID, _ := cmd.Flags().GetString("chat-id")
			text, _ := cmd.Flags().GetString("text")
			msgType, _ := cmd.Flags().GetString("type")
			sourceTarget, _ := cmd.Flags().GetString("source")
			body := map[string]interface{}{
				"chat_id": chatID,
				"text":    text,
				"type":    msgType,
				"source":  map[string]interface{}{"target": sourceTarget},
			}
			item, err := c.Post(context.Background(), "/chatmessages", body)
			if err != nil {
				return fmt.Errorf("could not create chatmessage: %w", err)
			}
			return output.PrintItem(cmd, item, chatmessageDetailColumns)
		},
	}
	cmd.Flags().String("chat-id", "", "Chat ID")
	cmd.Flags().String("text", "", "Message text")
	cmd.Flags().String("type", "normal", "Message type")
	cmd.Flags().String("source", "", "Source address target")
	_ = cmd.MarkFlagRequired("chat-id")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}

func newChatmessagesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a chat message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/chatmessages/"+args[0]); err != nil {
				return fmt.Errorf("could not delete chatmessage: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chat message %s deleted.\n", args[0])
			return nil
		},
	}
}
