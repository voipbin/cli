package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
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
			items, nextToken, err := c.List(context.Background(), "/queues", params)
			if err != nil {
				return fmt.Errorf("could not list queues: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, queueListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/queues/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get queue: %w", err)
			}
			return output.PrintItem(cmd, item, queueDetailColumns)
		},
	}
}

func newQueuesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new queue",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			routingMethod, _ := cmd.Flags().GetString("routing-method")
			serviceTimeout, _ := cmd.Flags().GetInt("service-timeout")
			waitTimeout, _ := cmd.Flags().GetInt("wait-timeout")
			waitFlowID, _ := cmd.Flags().GetString("wait-flow-id")
			body := map[string]interface{}{
				"name":            name,
				"detail":          detail,
				"routing_method":  routingMethod,
				"service_timeout": serviceTimeout,
				"wait_timeout":    waitTimeout,
				"wait_flow_id":    waitFlowID,
				"tag_ids":         []string{},
			}
			item, err := c.Post(context.Background(), "/queues", body)
			if err != nil {
				return fmt.Errorf("could not create queue: %w", err)
			}
			return output.PrintItem(cmd, item, queueDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			routingMethod, _ := cmd.Flags().GetString("routing-method")
			serviceTimeout, _ := cmd.Flags().GetInt("service-timeout")
			waitTimeout, _ := cmd.Flags().GetInt("wait-timeout")
			waitFlowID, _ := cmd.Flags().GetString("wait-flow-id")
			tagIDsJSON, _ := cmd.Flags().GetString("tag-ids")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if routingMethod != "" {
				body["routing_method"] = routingMethod
			}
			if serviceTimeout != 0 {
				body["service_timeout"] = serviceTimeout
			}
			if waitTimeout != 0 {
				body["wait_timeout"] = waitTimeout
			}
			if waitFlowID != "" {
				body["wait_flow_id"] = waitFlowID
			}
			if tagIDsJSON != "" {
				var parsed []interface{}
				if err := json.Unmarshal([]byte(tagIDsJSON), &parsed); err != nil {
					return fmt.Errorf("invalid JSON for --tag-ids: %w", err)
				}
				body["tag_ids"] = parsed
			}

			item, err := c.Put(context.Background(), "/queues/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update queue: %w", err)
			}
			return output.PrintItem(cmd, item, queueDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("routing-method", "", "New routing method")
	cmd.Flags().Int("service-timeout", 0, "New service timeout in ms")
	cmd.Flags().Int("wait-timeout", 0, "New wait timeout in ms")
	cmd.Flags().String("wait-flow-id", "", "New wait flow ID")
	cmd.Flags().String("tag-ids", "", "Tag IDs as JSON array")
	return cmd
}

func newQueuesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/queues/"+args[0]); err != nil {
				return fmt.Errorf("could not delete queue: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			routingMethod, _ := cmd.Flags().GetString("routing-method")
			body := map[string]interface{}{
				"routing_method": routingMethod,
			}
			if _, err := c.Put(context.Background(), "/queues/"+args[0]+"/routing_method", body); err != nil {
				return fmt.Errorf("could not set routing method: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			tagIDs, _ := cmd.Flags().GetStringSlice("tag-ids")
			body := map[string]interface{}{
				"tag_ids": tagIDs,
			}
			if _, err := c.Put(context.Background(), "/queues/"+args[0]+"/tag_ids", body); err != nil {
				return fmt.Errorf("could not set tag IDs: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Queue %s tag IDs set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().StringSlice("tag-ids", []string{}, "Comma-separated tag IDs")
	_ = cmd.MarkFlagRequired("tag-ids")
	return cmd
}
