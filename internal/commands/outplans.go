package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetOutplansParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetOutplansWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list outplans: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, outplanListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetOutplansIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get outplan: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outplanDetailColumns)
		},
	}
}

func newOutplansCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new outplan",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			dialTimeout, _ := cmd.Flags().GetInt("dial-timeout")
			tryInterval, _ := cmd.Flags().GetInt("try-interval")
			sourceTarget, _ := cmd.Flags().GetString("source")
			body := voipbin_client.PostOutplansJSONRequestBody{
				Name:        name,
				Detail:      detail,
				DialTimeout: dialTimeout,
				TryInterval: tryInterval,
				Source:      voipbin_client.CommonAddress{Target: &sourceTarget},
			}
			resp, err := client.PostOutplansWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create outplan: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outplanDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			body := voipbin_client.PutOutplansIdJSONRequestBody{
				Name:   name,
				Detail: detail,
			}
			resp, err := client.PutOutplansIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update outplan: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outplanDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteOutplansIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete outplan: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
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
			body := voipbin_client.PutOutplansIdDialInfoJSONRequestBody{
				DialTimeout:  dialTimeout,
				TryInterval:  tryInterval,
				Source:       voipbin_client.CommonAddress{Target: &sourceTarget},
				MaxTryCount0: maxTry0,
				MaxTryCount1: maxTry1,
				MaxTryCount2: maxTry2,
				MaxTryCount3: maxTry3,
				MaxTryCount4: maxTry4,
			}
			resp, err := client.PutOutplansIdDialInfoWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set dial info: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outplanDetailColumns)
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
