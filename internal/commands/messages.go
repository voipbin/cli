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

			items, nextToken, err := c.List(context.Background(), "/messages", params)
			if err != nil {
				return fmt.Errorf("could not list messages: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, messageListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/messages/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get message: %w", err)
			}

			return output.PrintItem(cmd, result, messageDetailColumns)
		},
	}
}

func newMessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new message",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			source, _ := cmd.Flags().GetString("source")
			destination, _ := cmd.Flags().GetString("destination")
			text, _ := cmd.Flags().GetString("text")

			body := map[string]interface{}{
				"source":       map[string]interface{}{"target": source},
				"destinations": []map[string]interface{}{{"target": destination}},
				"text":         text,
			}

			result, err := c.Post(context.Background(), "/messages", body)
			if err != nil {
				return fmt.Errorf("could not create message: %w", err)
			}

			return output.PrintItem(cmd, result, messageDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/messages/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete message: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Message %s deleted.\n", args[0])
			return nil
		},
	}
}
