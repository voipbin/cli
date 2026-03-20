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

func newRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Manage routes",
	}
	cmd.AddCommand(
		newRoutesListCmd(),
		newRoutesGetCmd(),
		newRoutesCreateCmd(),
		newRoutesUpdateCmd(),
		newRoutesDeleteCmd(),
	)
	return cmd
}

var routeListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "TARGET", Field: "target"},
	{Name: "PROVIDER_ID", Field: "provider_id"},
	{Name: "PRIORITY", Field: "priority"},
	{Name: "CREATED", Field: "tm_create"},
}

var routeDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "TARGET", Field: "target"},
	{Name: "PROVIDER_ID", Field: "provider_id"},
	{Name: "PRIORITY", Field: "priority"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newRoutesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			customerID, _ := cmd.Flags().GetString("customer-id")
			params := url.Values{}
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}
			if customerID != "" {
				params.Set("customer_id", customerID)
			}
			items, nextToken, err := c.List(context.Background(), "/routes", params)
			if err != nil {
				return fmt.Errorf("could not list routes: %w", err)
			}
			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}
			return output.PrintList(cmd, items, routeListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	cmd.Flags().String("customer-id", "", "Filter by customer ID")
	return cmd
}

func newRoutesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a route by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			item, err := c.Get(context.Background(), "/routes/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get route: %w", err)
			}
			return output.PrintItem(cmd, item, routeDetailColumns)
		},
	}
}

func newRoutesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new route",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			target, _ := cmd.Flags().GetString("target")
			providerID, _ := cmd.Flags().GetString("provider-id")
			priority, _ := cmd.Flags().GetInt("priority")
			customerID, _ := cmd.Flags().GetString("customer-id")
			body := map[string]interface{}{
				"name":        name,
				"detail":      detail,
				"target":      target,
				"provider_id": providerID,
				"priority":    priority,
				"customer_id": customerID,
			}
			item, err := c.Post(context.Background(), "/routes", body)
			if err != nil {
				return fmt.Errorf("could not create route: %w", err)
			}
			return output.PrintItem(cmd, item, routeDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Route name")
	cmd.Flags().String("detail", "", "Route detail")
	cmd.Flags().String("target", "", "Route target (e.g., country code or 'all')")
	cmd.Flags().String("provider-id", "", "Provider ID")
	cmd.Flags().Int("priority", 0, "Route priority")
	cmd.Flags().String("customer-id", "", "Customer ID")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("target")
	_ = cmd.MarkFlagRequired("provider-id")
	return cmd
}

func newRoutesUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a route",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			target, _ := cmd.Flags().GetString("target")
			providerID, _ := cmd.Flags().GetString("provider-id")
			priority, _ := cmd.Flags().GetInt("priority")
			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if target != "" {
				body["target"] = target
			}
			if providerID != "" {
				body["provider_id"] = providerID
			}
			if cmd.Flags().Changed("priority") {
				body["priority"] = priority
			}
			item, err := c.Put(context.Background(), "/routes/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update route: %w", err)
			}
			return output.PrintItem(cmd, item, routeDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	cmd.Flags().String("target", "", "New target")
	cmd.Flags().String("provider-id", "", "New provider ID")
	cmd.Flags().Int("priority", 0, "New priority")
	return cmd
}

func newRoutesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a route",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			if _, err := c.Delete(context.Background(), "/routes/"+args[0]); err != nil {
				return fmt.Errorf("could not delete route: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Route %s deleted.\n", args[0])
			return nil
		},
	}
}
