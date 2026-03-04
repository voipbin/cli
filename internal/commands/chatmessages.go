package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			chatID, _ := cmd.Flags().GetString("chat-id")
			params := voipbin_client.GetChatmessagesParams{
				ChatId: chatID,
			}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetChatmessagesWithResponse(context.Background(), &params)
			if err != nil {
				return fmt.Errorf("could not list chatmessages: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, chatmessageListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetChatmessagesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get chatmessage: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatmessageDetailColumns)
		},
	}
}

func newChatmessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chat message",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			chatID, _ := cmd.Flags().GetString("chat-id")
			text, _ := cmd.Flags().GetString("text")
			msgType, _ := cmd.Flags().GetString("type")
			sourceTarget, _ := cmd.Flags().GetString("source")
			body := voipbin_client.PostChatmessagesJSONRequestBody{
				ChatId: chatID,
				Text:   text,
				Type:   voipbin_client.ChatManagerMessagechatType(msgType),
				Source: voipbin_client.CommonAddress{Target: &sourceTarget},
			}
			resp, err := client.PostChatmessagesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create chatmessage: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatmessageDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteChatmessagesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete chatmessage: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chat message %s deleted.\n", args[0])
			return nil
		},
	}
}
