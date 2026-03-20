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

func newGroupcallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "groupcalls",
		Short: "Manage group calls",
	}
	cmd.AddCommand(
		newGroupcallsListCmd(),
		newGroupcallsGetCmd(),
		newGroupcallsCreateCmd(),
		newGroupcallsDeleteCmd(),
		newGroupcallsHangupCmd(),
	)
	return cmd
}

var groupcallListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "STATUS", Field: "status"},
	{Name: "RING_METHOD", Field: "ring_method"},
	{Name: "ANSWER_METHOD", Field: "answer_method"},
	{Name: "FLOW_ID", Field: "flow_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var groupcallDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "STATUS", Field: "status"},
	{Name: "RING_METHOD", Field: "ring_method"},
	{Name: "ANSWER_METHOD", Field: "answer_method"},
	{Name: "FLOW_ID", Field: "flow_id"},
	{Name: "ANSWER_CALL_ID", Field: "answer_call_id"},
	{Name: "MASTER_CALL_ID", Field: "master_call_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newGroupcallsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List group calls",
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

			items, nextToken, err := c.List(context.Background(), "/groupcalls", params)
			if err != nil {
				return fmt.Errorf("could not list group calls: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, groupcallListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newGroupcallsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a group call by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/groupcalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get group call: %w", err)
			}

			return output.PrintItem(cmd, result, groupcallDetailColumns)
		},
	}
}

func newGroupcallsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new group call",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			source, _ := cmd.Flags().GetString("source")
			destination, _ := cmd.Flags().GetString("destination")
			flowID, _ := cmd.Flags().GetString("flow-id")
			ringMethod, _ := cmd.Flags().GetString("ring-method")
			answerMethod, _ := cmd.Flags().GetString("answer-method")

			body := map[string]interface{}{
				"source":        map[string]interface{}{"target": source},
				"destinations":  []map[string]interface{}{{"target": destination}},
				"flow_id":       flowID,
				"ring_method":   ringMethod,
				"answer_method": answerMethod,
				"actions":       []interface{}{},
			}

			result, err := c.Post(context.Background(), "/groupcalls", body)
			if err != nil {
				return fmt.Errorf("could not create group call: %w", err)
			}

			return output.PrintItem(cmd, result, groupcallDetailColumns)
		},
	}
	cmd.Flags().String("source", "", "Source address")
	cmd.Flags().String("destination", "", "Destination address")
	cmd.Flags().String("flow-id", "", "Flow ID")
	cmd.Flags().String("ring-method", "", "Ring method")
	cmd.Flags().String("answer-method", "", "Answer method")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("destination")
	return cmd
}

func newGroupcallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a group call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/groupcalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete group call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Group call %s deleted.\n", args[0])
			return nil
		},
	}
}

func newGroupcallsHangupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hangup <id>",
		Short: "Hang up a group call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/groupcalls/"+args[0]+"/hangup", nil)
			if err != nil {
				return fmt.Errorf("could not hangup group call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Group call %s hung up.\n", args[0])
			return nil
		},
	}
}
