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

func newChatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chats",
		Short: "Manage chats",
	}
	cmd.AddCommand(
		newChatsListCmd(),
		newChatsGetCmd(),
		newChatsCreateCmd(),
		newChatsUpdateCmd(),
		newChatsDeleteCmd(),
		newChatsAddParticipantCmd(),
		newChatsRemoveParticipantCmd(),
		newChatsSetRoomOwnerCmd(),
	)
	return cmd
}

var chatListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "TYPE", Field: "type"},
	{Name: "ROOM_OWNER_ID", Field: "room_owner_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var chatDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "TYPE", Field: "type"},
	{Name: "ROOM_OWNER_ID", Field: "room_owner_id"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newChatsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List chats",
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
			items, nextToken, err := c.List(context.Background(), "/chats", params)
			if err != nil {
				return fmt.Errorf("could not list chats: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, chatListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newChatsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a chat by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/chats/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get chat: %w", err)
			}
			return output.PrintItem(cmd, item, chatDetailColumns)
		},
	}
}

func newChatsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			ownerID, _ := cmd.Flags().GetString("owner-id")
			chatType, _ := cmd.Flags().GetString("type")
			participantIDs, _ := cmd.Flags().GetStringSlice("participant-ids")
			body := map[string]interface{}{
				"name":            name,
				"detail":          detail,
				"owner_id":        ownerID,
				"type":            chatType,
				"participant_ids": participantIDs,
			}
			item, err := c.Post(context.Background(), "/chats", body)
			if err != nil {
				return fmt.Errorf("could not create chat: %w", err)
			}
			return output.PrintItem(cmd, item, chatDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Chat name")
	cmd.Flags().String("detail", "", "Chat detail")
	cmd.Flags().String("owner-id", "", "Owner agent ID")
	cmd.Flags().String("type", "direct", "Chat type")
	cmd.Flags().StringSlice("participant-ids", []string{}, "Participant agent IDs")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("owner-id")
	return cmd
}

func newChatsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a chat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			body := map[string]interface{}{
				"name":   name,
				"detail": detail,
			}
			item, err := c.Put(context.Background(), "/chats/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update chat: %w", err)
			}
			return output.PrintItem(cmd, item, chatDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	return cmd
}

func newChatsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a chat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/chats/"+args[0]); err != nil {
				return fmt.Errorf("could not delete chat: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chat %s deleted.\n", args[0])
			return nil
		},
	}
}

func newChatsAddParticipantCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-participant <id>",
		Short: "Add a participant to a chat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			participantID, _ := cmd.Flags().GetString("participant-id")
			body := map[string]interface{}{
				"participant_id": participantID,
			}
			if _, err := c.Post(context.Background(), "/chats/"+args[0]+"/participant_ids", body); err != nil {
				return fmt.Errorf("could not add participant: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Participant %s added to chat %s.\n", participantID, args[0])
			return nil
		},
	}
	cmd.Flags().String("participant-id", "", "Participant agent ID")
	_ = cmd.MarkFlagRequired("participant-id")
	return cmd
}

func newChatsRemoveParticipantCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-participant <id> <participant-id>",
		Short: "Remove a participant from a chat",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/chats/"+args[0]+"/participant_ids/"+args[1]); err != nil {
				return fmt.Errorf("could not remove participant: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Participant %s removed from chat %s.\n", args[1], args[0])
			return nil
		},
	}
}

func newChatsSetRoomOwnerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-room-owner <id>",
		Short: "Set room owner for a chat",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			roomOwnerID, _ := cmd.Flags().GetString("room-owner-id")
			body := map[string]interface{}{
				"room_owner_id": roomOwnerID,
			}
			if _, err := c.Put(context.Background(), "/chats/"+args[0]+"/room_owner_id", body); err != nil {
				return fmt.Errorf("could not set room owner: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chat %s room owner set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("room-owner-id", "", "Room owner agent ID")
	_ = cmd.MarkFlagRequired("room-owner-id")
	return cmd
}
