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

func newEmailsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "emails",
		Short: "Manage emails",
	}
	cmd.AddCommand(
		newEmailsListCmd(),
		newEmailsGetCmd(),
		newEmailsCreateCmd(),
		newEmailsDeleteCmd(),
	)
	return cmd
}

var emailListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "SUBJECT", Field: "subject"},
	{Name: "STATUS", Field: "status"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var emailDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "SUBJECT", Field: "subject"},
	{Name: "CONTENT", Field: "content"},
	{Name: "STATUS", Field: "status"},
	{Name: "ACTIVEFLOW_ID", Field: "activeflow_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newEmailsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List emails",
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

			items, nextToken, err := c.List(context.Background(), "/emails", params)
			if err != nil {
				return fmt.Errorf("could not list emails: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, emailListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newEmailsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an email by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/emails/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get email: %w", err)
			}

			return output.PrintItem(cmd, result, emailDetailColumns)
		},
	}
}

func newEmailsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new email",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			destination, _ := cmd.Flags().GetString("destination")
			subject, _ := cmd.Flags().GetString("subject")
			content, _ := cmd.Flags().GetString("content")

			body := map[string]interface{}{
				"destinations": []map[string]interface{}{{"target": destination}},
				"subject":      subject,
				"content":      content,
				"attachments":  []interface{}{},
			}

			result, err := c.Post(context.Background(), "/emails", body)
			if err != nil {
				return fmt.Errorf("could not create email: %w", err)
			}

			return output.PrintItem(cmd, result, emailDetailColumns)
		},
	}
	cmd.Flags().String("destination", "", "Destination email address")
	cmd.Flags().String("subject", "", "Email subject")
	cmd.Flags().String("content", "", "Email content")
	_ = cmd.MarkFlagRequired("destination")
	_ = cmd.MarkFlagRequired("subject")
	_ = cmd.MarkFlagRequired("content")
	return cmd
}

func newEmailsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an email",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/emails/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete email: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Email %s deleted.\n", args[0])
			return nil
		},
	}
}
