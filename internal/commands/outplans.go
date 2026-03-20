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

func newOutplansCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outplans",
		Short: "Manage outplans",
	}
	cmd.AddCommand(
		newOutplansListCmd(),
		newOutplansGetCmd(),
		newOutplansCreateCmd(),
		newOutplansUpdateCmd(),
		newOutplansDeleteCmd(),
		newOutplansSetDialInfoCmd(),
	)
	return cmd
}

var outplanListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DIAL_TIMEOUT", Field: "dial_timeout"},
	{Name: "TRY_INTERVAL", Field: "try_interval"},
	{Name: "CREATED", Field: "tm_create"},
}

var outplanDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "DIAL_TIMEOUT", Field: "dial_timeout"},
	{Name: "TRY_INTERVAL", Field: "try_interval"},
	{Name: "MAX_TRY_COUNT_0", Field: "max_try_count_0"},
	{Name: "MAX_TRY_COUNT_1", Field: "max_try_count_1"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newOutplansListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List outplans",
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
			items, nextToken, err := c.List(context.Background(), "/outplans", params)
			if err != nil {
				return fmt.Errorf("could not list outplans: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, outplanListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newOutplansGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an outplan by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			result, err := c.Get(context.Background(), "/outplans/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get outplan: %w", err)
			}
			return output.PrintItem(cmd, result, outplanDetailColumns)
		},
	}
}

func newOutplansCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new outplan",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			dialTimeout, _ := cmd.Flags().GetInt("dial-timeout")
			tryInterval, _ := cmd.Flags().GetInt("try-interval")
			sourceTarget, _ := cmd.Flags().GetString("source")
			body := map[string]interface{}{
				"name":         name,
				"detail":       detail,
				"dial_timeout": dialTimeout,
				"try_interval": tryInterval,
				"source":       map[string]interface{}{"target": sourceTarget},
			}
			result, err := c.Post(context.Background(), "/outplans", body)
			if err != nil {
				return fmt.Errorf("could not create outplan: %w", err)
			}
			return output.PrintItem(cmd, result, outplanDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Outplan name")
	cmd.Flags().String("detail", "", "Outplan detail")
	cmd.Flags().Int("dial-timeout", 30, "Dial timeout in seconds")
	cmd.Flags().Int("try-interval", 60, "Interval between retry attempts")
	cmd.Flags().String("source", "", "Source address target")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newOutplansUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an outplan",
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
			result, err := c.Put(context.Background(), "/outplans/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update outplan: %w", err)
			}
			return output.PrintItem(cmd, result, outplanDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	return cmd
}

func newOutplansDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an outplan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			_, err = c.Delete(context.Background(), "/outplans/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete outplan: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Outplan %s deleted.\n", args[0])
			return nil
		},
	}
}

func newOutplansSetDialInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-dial-info <id>",
		Short: "Update dial info for an outplan",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			dialTimeout, _ := cmd.Flags().GetInt("dial-timeout")
			tryInterval, _ := cmd.Flags().GetInt("try-interval")
			sourceTarget, _ := cmd.Flags().GetString("source")
			maxTry0, _ := cmd.Flags().GetInt("max-try-0")
			maxTry1, _ := cmd.Flags().GetInt("max-try-1")
			maxTry2, _ := cmd.Flags().GetInt("max-try-2")
			maxTry3, _ := cmd.Flags().GetInt("max-try-3")
			maxTry4, _ := cmd.Flags().GetInt("max-try-4")
			body := map[string]interface{}{}
			if cmd.Flags().Changed("dial-timeout") {
				body["dial_timeout"] = dialTimeout
			}
			if cmd.Flags().Changed("try-interval") {
				body["try_interval"] = tryInterval
			}
			if sourceTarget != "" {
				body["source"] = map[string]interface{}{"target": sourceTarget}
			}
			if cmd.Flags().Changed("max-try-0") {
				body["max_try_count0"] = maxTry0
			}
			if cmd.Flags().Changed("max-try-1") {
				body["max_try_count1"] = maxTry1
			}
			if cmd.Flags().Changed("max-try-2") {
				body["max_try_count2"] = maxTry2
			}
			if cmd.Flags().Changed("max-try-3") {
				body["max_try_count3"] = maxTry3
			}
			if cmd.Flags().Changed("max-try-4") {
				body["max_try_count4"] = maxTry4
			}
			result, err := c.Put(context.Background(), "/outplans/"+args[0]+"/dial_info", body)
			if err != nil {
				return fmt.Errorf("could not set dial info: %w", err)
			}
			return output.PrintItem(cmd, result, outplanDetailColumns)
		},
	}
	cmd.Flags().Int("dial-timeout", 30, "Dial timeout in seconds")
	cmd.Flags().Int("try-interval", 60, "Interval between retries")
	cmd.Flags().String("source", "", "Source address target")
	cmd.Flags().Int("max-try-0", 1, "Max try count for destination 0")
	cmd.Flags().Int("max-try-1", 0, "Max try count for destination 1")
	cmd.Flags().Int("max-try-2", 0, "Max try count for destination 2")
	cmd.Flags().Int("max-try-3", 0, "Max try count for destination 3")
	cmd.Flags().Int("max-try-4", 0, "Max try count for destination 4")
	return cmd
}
