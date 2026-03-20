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

func newAimessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aimessages",
		Short: "Manage AI messages",
	}
	cmd.AddCommand(
		newAimessagesListCmd(),
		newAimessagesGetCmd(),
		newAimessagesCreateCmd(),
		newAimessagesDeleteCmd(),
	)
	return cmd
}

var aimessageListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "AICALL_ID", Field: "aicall_id"},
	{Name: "ROLE", Field: "role"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "CONTENT", Field: "content"},
	{Name: "CREATED", Field: "tm_create"},
}

var aimessageDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "AICALL_ID", Field: "aicall_id"},
	{Name: "ROLE", Field: "role"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "CONTENT", Field: "content"},
	{Name: "CREATED", Field: "tm_create"},
}

func newAimessagesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List AI messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			aicallID, _ := cmd.Flags().GetString("aicall-id")

			params := url.Values{}
			params.Set("aicall_id", aicallID)
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}

			items, nextToken, err := c.List(context.Background(), "/aimessages", params)
			if err != nil {
				return fmt.Errorf("could not list AI messages: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, aimessageListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	cmd.Flags().String("aicall-id", "", "AI call ID to filter by")
	_ = cmd.MarkFlagRequired("aicall-id")
	return cmd
}

func newAimessagesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an AI message by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/aimessages/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get AI message: %w", err)
			}

			return output.PrintItem(cmd, result, aimessageDetailColumns)
		},
	}
}

func newAimessagesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI message",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			aicallID, _ := cmd.Flags().GetString("aicall-id")
			content, _ := cmd.Flags().GetString("content")
			role, _ := cmd.Flags().GetString("role")

			body := map[string]interface{}{
				"aicall_id": aicallID,
				"content":   content,
				"role":      role,
			}

			result, err := c.Post(context.Background(), "/aimessages", body)
			if err != nil {
				return fmt.Errorf("could not create AI message: %w", err)
			}

			return output.PrintItem(cmd, result, aimessageDetailColumns)
		},
	}
	cmd.Flags().String("aicall-id", "", "AI call ID")
	cmd.Flags().String("content", "", "Message content")
	cmd.Flags().String("role", "", "Message role")
	_ = cmd.MarkFlagRequired("aicall-id")
	_ = cmd.MarkFlagRequired("content")
	_ = cmd.MarkFlagRequired("role")
	return cmd
}

func newAimessagesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an AI message",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/aimessages/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI message: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI message %s deleted.\n", args[0])
			return nil
		},
	}
}
