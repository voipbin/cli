package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newChatroomsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chatrooms",
		Short: "Manage chatrooms",
	}
	cmd.AddCommand(
		newChatroomsListCmd(),
		newChatroomsGetCmd(),
		newChatroomsCreateCmd(),
		newChatroomsUpdateCmd(),
		newChatroomsDeleteCmd(),
	)
	return cmd
}

var chatroomListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "CHAT_ID", Field: "chat_id"},
	{Name: "ROOM_OWNER_ID", Field: "room_owner_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var chatroomDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CHAT_ID", Field: "chat_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "ROOM_OWNER_ID", Field: "room_owner_id"},
	{Name: "TYPE", Field: "type"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newChatroomsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List chatrooms",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			ownerID, _ := cmd.Flags().GetString("owner-id")
			params := &voipbin_client.GetChatroomsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			if ownerID != "" {
				params.OwnerId = &ownerID
			}
			resp, err := client.GetChatroomsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list chatrooms: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, chatroomListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	cmd.Flags().String("owner-id", "", "Filter by owner ID")
	return cmd
}

func newChatroomsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a chatroom by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetChatroomsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get chatroom: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatroomDetailColumns)
		},
	}
}

func newChatroomsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chatroom",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			participantIDs, _ := cmd.Flags().GetStringSlice("participant-ids")
			body := voipbin_client.PostChatroomsJSONRequestBody{
				Name:           name,
				Detail:         detail,
				ParticipantIds: participantIDs,
			}
			resp, err := client.PostChatroomsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create chatroom: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatroomDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Chatroom name")
	cmd.Flags().String("detail", "", "Chatroom detail")
	cmd.Flags().StringSlice("participant-ids", []string{}, "Participant agent IDs")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newChatroomsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a chatroom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			body := voipbin_client.PutChatroomsIdJSONRequestBody{
				Name:   name,
				Detail: detail,
			}
			resp, err := client.PutChatroomsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update chatroom: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, chatroomDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	return cmd
}

func newChatroomsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a chatroom",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteChatroomsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete chatroom: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chatroom %s deleted.\n", args[0])
			return nil
		},
	}
}
