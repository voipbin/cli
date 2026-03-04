package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			chatroomID, _ := cmd.Flags().GetString("chatroom-id")
			params := voipbin_client.GetChatroommessagesParams{
				ChatroomId: chatroomID,
			}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetChatroommessagesWithResponse(context.Background(), &params)
			if err != nil {
				return fmt.Errorf("could not list chatroommessages: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, chatroommessageListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetChatroommessagesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get chatroommessage: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatroommessageDetailColumns)
		},
	}
}

func newChatroommessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chatroom message",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			chatroomID, _ := cmd.Flags().GetString("chatroom-id")
			text, _ := cmd.Flags().GetString("text")
			body := voipbin_client.PostChatroommessagesJSONRequestBody{
				ChatroomId: chatroomID,
				Text:       text,
			}
			resp, err := client.PostChatroommessagesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create chatroommessage: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatroommessageDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteChatroommessagesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete chatroommessage: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chatroom message %s deleted.\n", args[0])
			return nil
		},
	}
}
