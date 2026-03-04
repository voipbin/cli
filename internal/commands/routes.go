package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			customerID, _ := cmd.Flags().GetString("customer-id")
			params := &voipbin_client.GetRoutesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			if customerID != "" {
				params.CustomerId = &customerID
			}
			resp, err := client.GetRoutesWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list routes: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, routeListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetRoutesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get route: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, routeDetailColumns)
		},
	}
}

func newRoutesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new route",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			target, _ := cmd.Flags().GetString("target")
			providerID, _ := cmd.Flags().GetString("provider-id")
			priority, _ := cmd.Flags().GetInt("priority")
			customerID, _ := cmd.Flags().GetString("customer-id")
			body := voipbin_client.PostRoutesJSONRequestBody{
				Name:       name,
				Detail:     detail,
				Target:     target,
				ProviderId: providerID,
				Priority:   priority,
				CustomerId: customerID,
			}
			resp, err := client.PostRoutesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create route: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, routeDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			target, _ := cmd.Flags().GetString("target")
			providerID, _ := cmd.Flags().GetString("provider-id")
			priority, _ := cmd.Flags().GetInt("priority")
			body := voipbin_client.PutRoutesIdJSONRequestBody{
				Name:       name,
				Detail:     detail,
				Target:     target,
				ProviderId: providerID,
				Priority:   priority,
			}
			resp, err := client.PutRoutesIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update route: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, routeDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteRoutesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete route: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Route %s deleted.\n", args[0])
			return nil
		},
	}
}
