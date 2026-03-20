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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			ownerID, _ := cmd.Flags().GetString("owner-id")
			params := url.Values{}
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}
			if ownerID != "" {
				params.Set("owner_id", ownerID)
			}
			items, nextToken, err := c.List(context.Background(), "/chatrooms", params)
			if err != nil {
				return fmt.Errorf("could not list chatrooms: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, chatroomListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/chatrooms/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get chatroom: %w", err)
			}
			return output.PrintItem(cmd, item, chatroomDetailColumns)
		},
	}
}

func newChatroomsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new chatroom",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			participantIDs, _ := cmd.Flags().GetStringSlice("participant-ids")
			body := map[string]interface{}{
				"name":            name,
				"detail":          detail,
				"participant_ids": participantIDs,
			}
			item, err := c.Post(context.Background(), "/chatrooms", body)
			if err != nil {
				return fmt.Errorf("could not create chatroom: %w", err)
			}
			return output.PrintItem(cmd, item, chatroomDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			item, err := c.Put(context.Background(), "/chatrooms/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update chatroom: %w", err)
			}
			return output.PrintItem(cmd, item, chatroomDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/chatrooms/"+args[0]); err != nil {
				return fmt.Errorf("could not delete chatroom: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Chatroom %s deleted.\n", args[0])
			return nil
		},
	}
}
