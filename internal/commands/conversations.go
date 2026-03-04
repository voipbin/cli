package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetConversationsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			body := voipbin_client.GetConversationsJSONRequestBody{}
			resp, err := client.GetConversationsWithResponse(context.Background(), params, body)
			if err != nil {
				return fmt.Errorf("could not list conversations: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, conversationListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetConversationsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get conversation: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, conversationDetailColumns)
		},
	}
}

func newConversationsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a conversation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			ownerID, _ := cmd.Flags().GetString("owner-id")
			ownerType, _ := cmd.Flags().GetString("owner-type")
			body := voipbin_client.PutConversationsIdJSONRequestBody{}
			if name != "" {
				body.Name = &name
			}
			if detail != "" {
				body.Detail = &detail
			}
			if ownerID != "" {
				body.OwnerId = &ownerID
			}
			if ownerType != "" {
				body.OwnerType = &ownerType
			}
			resp, err := client.PutConversationsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update conversation: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, conversationDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetConversationsIdMessagesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetConversationsIdMessagesWithResponse(context.Background(), args[0], params)
			if err != nil {
				return fmt.Errorf("could not list messages: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, conversationMessageColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			text, _ := cmd.Flags().GetString("text")
			body := voipbin_client.PostConversationsIdMessagesJSONRequestBody{
				Text:   text,
				Medias: []voipbin_client.ConversationManagerMedia{},
			}
			resp, err := client.PostConversationsIdMessagesWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not create message: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, conversationMessageColumns)
		},
	}
	cmd.Flags().String("text", "", "Message text")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}
