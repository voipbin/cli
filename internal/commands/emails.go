package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetEmailsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetEmailsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list emails: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, emailListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetEmailsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get email: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, emailDetailColumns)
		},
	}
}

func newEmailsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new email",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			destination, _ := cmd.Flags().GetString("destination")
			subject, _ := cmd.Flags().GetString("subject")
			content, _ := cmd.Flags().GetString("content")

			body := voipbin_client.PostEmailsJSONRequestBody{
				Destinations: []voipbin_client.CommonAddress{{Target: &destination}},
				Subject:      subject,
				Content:      content,
				Attachments:  []voipbin_client.EmailManagerEmailAttachment{},
			}

			resp, err := client.PostEmailsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create email: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, emailDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteEmailsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete email: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Email %s deleted.\n", args[0])
			return nil
		},
	}
}
