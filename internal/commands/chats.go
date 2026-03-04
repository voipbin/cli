package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetChatsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetChatsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list chats: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, chatListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetChatsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get chat: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatDetailColumns)
		},
	}
}

func newChatsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chat",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			ownerID, _ := cmd.Flags().GetString("owner-id")
			chatType, _ := cmd.Flags().GetString("type")
			participantIDs, _ := cmd.Flags().GetStringSlice("participant-ids")
			body := voipbin_client.PostChatsJSONRequestBody{
				Name:           name,
				Detail:         detail,
				OwnerId:        ownerID,
				Type:           voipbin_client.ChatManagerChatType(chatType),
				ParticipantIds: participantIDs,
			}
			resp, err := client.PostChatsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create chat: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			body := voipbin_client.PutChatsIdJSONRequestBody{
				Name:   name,
				Detail: detail,
			}
			resp, err := client.PutChatsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update chat: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteChatsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete chat: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			participantID, _ := cmd.Flags().GetString("participant-id")
			body := voipbin_client.PostChatsIdParticipantIdsJSONRequestBody{
				ParticipantId: participantID,
			}
			resp, err := client.PostChatsIdParticipantIdsWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not add participant: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteChatsIdParticipantIdsParticipantIdWithResponse(context.Background(), args[0], args[1])
			if err != nil {
				return fmt.Errorf("could not remove participant: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			roomOwnerID, _ := cmd.Flags().GetString("room-owner-id")
			body := voipbin_client.PutChatsIdRoomOwnerIdJSONRequestBody{
				RoomOwnerId: roomOwnerID,
			}
			resp, err := client.PutChatsIdRoomOwnerIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set room owner: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chat %s room owner set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("room-owner-id", "", "Room owner agent ID")
	_ = cmd.MarkFlagRequired("room-owner-id")
	return cmd
}
