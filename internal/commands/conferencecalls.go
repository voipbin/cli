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

func newConferencecallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conferencecalls",
		Short: "Manage conference calls",
	}
	cmd.AddCommand(
		newConferencecallsListCmd(),
		newConferencecallsGetCmd(),
		newConferencecallsDeleteCmd(),
	)
	return cmd
}

var conferencecallListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CONFERENCE_ID", Field: "conference_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var conferencecallDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "CONFERENCE_ID", Field: "conference_id"},
	{Name: "ACTIVEFLOW_ID", Field: "activeflow_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

func newConferencecallsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List conference calls",
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

			items, nextToken, err := c.List(context.Background(), "/conferencecalls", params)
			if err != nil {
				return fmt.Errorf("could not list conference calls: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, conferencecallListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newConferencecallsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a conference call by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/conferencecalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get conference call: %w", err)
			}

			return output.PrintItem(cmd, result, conferencecallDetailColumns)
		},
	}
}

func newConferencecallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a conference call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/conferencecalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete conference call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conference call %s deleted.\n", args[0])
			return nil
		},
	}
}
