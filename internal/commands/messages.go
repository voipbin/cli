package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage messages",
	}
	cmd.AddCommand(
		newMessagesListCmd(),
		newMessagesGetCmd(),
		newMessagesCreateCmd(),
		newMessagesDeleteCmd(),
	)
	return cmd
}

var messageListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "TYPE", Field: "type"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
}

var messageDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "TYPE", Field: "type"},
	{Name: "TEXT", Field: "text"},
	{Name: "CREATED", Field: "tm_create"},
}

func newMessagesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetMessagesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetMessagesWithResponse(context.Background(), params)
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

			return output.PrintList(cmd, *resp.JSON200.Result, messageListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newMessagesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a message by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetMessagesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get message: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, messageDetailColumns)
		},
	}
}

func newMessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new message",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			source, _ := cmd.Flags().GetString("source")
			destination, _ := cmd.Flags().GetString("destination")
			text, _ := cmd.Flags().GetString("text")

			body := voipbin_client.PostMessagesJSONRequestBody{
				Source:       voipbin_client.CommonAddress{Target: &source},
				Destinations: []voipbin_client.CommonAddress{{Target: &destination}},
				Text:         text,
			}

			resp, err := client.PostMessagesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create message: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, messageDetailColumns)
		},
	}
	cmd.Flags().String("source", "", "Source address")
	cmd.Flags().String("destination", "", "Destination address")
	cmd.Flags().String("text", "", "Message text")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("destination")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}

func newMessagesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteMessagesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete message: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Message %s deleted.\n", args[0])
			return nil
		},
	}
}
