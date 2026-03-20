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

func newChatroommessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chatroommessages",
		Short: "Manage chatroom messages",
	}
	cmd.AddCommand(
		newChatroommessagesListCmd(),
		newChatroommessagesGetCmd(),
		newChatroommessagesCreateCmd(),
		newChatroommessagesDeleteCmd(),
	)
	return cmd
}

var chatroommessageListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CHATROOM_ID", Field: "chatroom_id"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
}

var chatroommessageDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CHATROOM_ID", Field: "chatroom_id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newChatroommessagesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List chatroom messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			chatroomID, _ := cmd.Flags().GetString("chatroom-id")
			params := url.Values{}
			params.Set("chatroom_id", chatroomID)
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}
			items, nextToken, err := c.List(context.Background(), "/chatroommessages", params)
			if err != nil {
				return fmt.Errorf("could not list chatroommessages: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, chatroommessageListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	cmd.Flags().String("chatroom-id", "", "Chatroom ID to list messages for")
	_ = cmd.MarkFlagRequired("chatroom-id")
	return cmd
}

func newChatroommessagesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a chatroom message by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/chatroommessages/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get chatroommessage: %w", err)
			}
			return output.PrintItem(cmd, item, chatroommessageDetailColumns)
		},
	}
}

func newChatroommessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chatroom message",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			chatroomID, _ := cmd.Flags().GetString("chatroom-id")
			text, _ := cmd.Flags().GetString("text")
			body := map[string]interface{}{
				"chatroom_id": chatroomID,
				"text":        text,
			}
			item, err := c.Post(context.Background(), "/chatroommessages", body)
			if err != nil {
				return fmt.Errorf("could not create chatroommessage: %w", err)
			}
			return output.PrintItem(cmd, item, chatroommessageDetailColumns)
		},
	}
	cmd.Flags().String("chatroom-id", "", "Chatroom ID")
	cmd.Flags().String("text", "", "Message text")
	_ = cmd.MarkFlagRequired("chatroom-id")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}

func newChatroommessagesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a chatroom message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/chatroommessages/"+args[0]); err != nil {
				return fmt.Errorf("could not delete chatroommessage: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chatroom message %s deleted.\n", args[0])
			return nil
		},
	}
}
