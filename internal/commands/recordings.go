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

func newRecordingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recordings",
		Short: "Manage recordings",
	}
	cmd.AddCommand(
		newRecordingsListCmd(),
		newRecordingsGetCmd(),
		newRecordingsDeleteCmd(),
	)
	return cmd
}

var recordingListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "OWNER_TYPE", Field: "owner_type"},
	{Name: "FORMAT", Field: "format"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var recordingDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "OWNER_TYPE", Field: "owner_type"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "FORMAT", Field: "format"},
	{Name: "STATUS", Field: "status"},
	{Name: "ACTIVEFLOW_ID", Field: "activeflow_id"},
	{Name: "ON_END_FLOW_ID", Field: "on_end_flow_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newRecordingsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recordings",
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

			items, nextToken, err := c.List(context.Background(), "/recordings", params)
			if err != nil {
				return fmt.Errorf("could not list recordings: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, recordingListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newRecordingsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a recording by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/recordings/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get recording: %w", err)
			}

			return output.PrintItem(cmd, result, recordingDetailColumns)
		},
	}
}

func newRecordingsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a recording",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/recordings/"+args[0]); err != nil {
				return fmt.Errorf("could not delete recording: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Recording %s deleted.\n", args[0])
			return nil
		},
	}
}
