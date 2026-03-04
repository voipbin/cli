package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newQueuesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queues",
		Short: "Manage queues",
	}
	cmd.AddCommand(
		newQueuesListCmd(),
		newQueuesGetCmd(),
		newQueuesCreateCmd(),
		newQueuesUpdateCmd(),
		newQueuesDeleteCmd(),
		newQueuesSetRoutingMethodCmd(),
		newQueuesSetTagIdsCmd(),
	)
	return cmd
}

var queueListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "ROUTING_METHOD", Field: "routing_method"},
	{Name: "SERVICE_TIMEOUT", Field: "service_timeout"},
	{Name: "CREATED", Field: "tm_create"},
}

var queueDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "ROUTING_METHOD", Field: "routing_method"},
	{Name: "SERVICE_TIMEOUT", Field: "service_timeout"},
	{Name: "WAIT_TIMEOUT", Field: "wait_timeout"},
	{Name: "WAIT_FLOW_ID", Field: "wait_flow_id"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newQueuesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List queues",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetQueuesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetQueuesWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list queues: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, queueListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newQueuesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a queue by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetQueuesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get queue: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, queueDetailColumns)
		},
	}
}

func newQueuesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			routingMethod, _ := cmd.Flags().GetString("routing-method")
			serviceTimeout, _ := cmd.Flags().GetInt("service-timeout")
			waitTimeout, _ := cmd.Flags().GetInt("wait-timeout")
			waitFlowID, _ := cmd.Flags().GetString("wait-flow-id")
			body := voipbin_client.PostQueuesJSONRequestBody{
				Name:           name,
				Detail:         detail,
				RoutingMethod:  voipbin_client.QueueManagerQueueRoutingMethod(routingMethod),
				ServiceTimeout: serviceTimeout,
				WaitTimeout:    waitTimeout,
				WaitFlowId:     waitFlowID,
				TagIds:         []string{},
			}
			resp, err := client.PostQueuesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create queue: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, queueDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Queue name")
	cmd.Flags().String("detail", "", "Queue detail")
	cmd.Flags().String("routing-method", "ringall", "Routing method")
	cmd.Flags().Int("service-timeout", 60000, "Service timeout in milliseconds")
	cmd.Flags().Int("wait-timeout", 300000, "Wait timeout in milliseconds")
	cmd.Flags().String("wait-flow-id", "", "Flow ID for wait queue")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newQueuesUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			routingMethod, _ := cmd.Flags().GetString("routing-method")
			serviceTimeout, _ := cmd.Flags().GetInt("service-timeout")
			waitTimeout, _ := cmd.Flags().GetInt("wait-timeout")
			waitFlowID, _ := cmd.Flags().GetString("wait-flow-id")
			body := voipbin_client.PutQueuesIdJSONRequestBody{
				Name:           name,
				Detail:         detail,
				RoutingMethod:  voipbin_client.QueueManagerQueueRoutingMethod(routingMethod),
				ServiceTimeout: serviceTimeout,
				WaitTimeout:    waitTimeout,
				WaitFlowId:     waitFlowID,
				TagIds:         []string{},
			}
			resp, err := client.PutQueuesIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update queue: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, queueDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("routing-method", "", "New routing method")
	cmd.Flags().Int("service-timeout", 0, "New service timeout in ms")
	cmd.Flags().Int("wait-timeout", 0, "New wait timeout in ms")
	cmd.Flags().String("wait-flow-id", "", "New wait flow ID")
	return cmd
}

func newQueuesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteQueuesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete queue: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queue %s deleted.\n", args[0])
			return nil
		},
	}
}

func newQueuesSetRoutingMethodCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-routing-method <id>",
		Short: "Set routing method for a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			routingMethod, _ := cmd.Flags().GetString("routing-method")
			body := voipbin_client.PutQueuesIdRoutingMethodJSONRequestBody{
				RoutingMethod: voipbin_client.QueueManagerQueueRoutingMethod(routingMethod),
			}
			resp, err := client.PutQueuesIdRoutingMethodWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set routing method: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queue %s routing method set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("routing-method", "", "Routing method")
	_ = cmd.MarkFlagRequired("routing-method")
	return cmd
}

func newQueuesSetTagIdsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-tag-ids <id>",
		Short: "Set tag IDs for a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			tagIDs, _ := cmd.Flags().GetStringSlice("tag-ids")
			body := voipbin_client.PutQueuesIdTagIdsJSONRequestBody{
				TagIds: tagIDs,
			}
			resp, err := client.PutQueuesIdTagIdsWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set tag IDs: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queue %s tag IDs set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringSlice("tag-ids", []string{}, "Comma-separated tag IDs")
	_ = cmd.MarkFlagRequired("tag-ids")
	return cmd
}
