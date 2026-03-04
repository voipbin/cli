package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newProvidersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "Manage providers",
	}
	cmd.AddCommand(
		newProvidersListCmd(),
		newProvidersGetCmd(),
		newProvidersCreateCmd(),
		newProvidersUpdateCmd(),
		newProvidersDeleteCmd(),
	)
	return cmd
}

var providerListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "HOSTNAME", Field: "hostname"},
	{Name: "TYPE", Field: "type"},
	{Name: "CREATED", Field: "tm_create"},
}

var providerDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "HOSTNAME", Field: "hostname"},
	{Name: "TYPE", Field: "type"},
	{Name: "TECH_PREFIX", Field: "tech_prefix"},
	{Name: "TECH_POSTFIX", Field: "tech_postfix"},
	{Name: "CREATED", Field: "tm_create"},
}

func newProvidersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetProvidersParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetProvidersWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list providers: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, providerListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newProvidersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a provider by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetProvidersIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get provider: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, providerDetailColumns)
		},
	}
}

func newProvidersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			hostname, _ := cmd.Flags().GetString("hostname")
			detail, _ := cmd.Flags().GetString("detail")
			techPrefix, _ := cmd.Flags().GetString("tech-prefix")
			techPostfix, _ := cmd.Flags().GetString("tech-postfix")
			provType, _ := cmd.Flags().GetString("type")

			body := voipbin_client.PostProvidersJSONRequestBody{
				Name:        name,
				Hostname:    hostname,
				Detail:      detail,
				TechPrefix:  techPrefix,
				TechPostfix: techPostfix,
				Type:        voipbin_client.RouteManagerProviderType(provType),
				TechHeaders: map[string]interface{}{},
			}

			resp, err := client.PostProvidersWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create provider: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, providerDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Provider name")
	cmd.Flags().String("hostname", "", "Destination hostname")
	cmd.Flags().String("detail", "", "Provider detail")
	cmd.Flags().String("tech-prefix", "", "Tech prefix")
	cmd.Flags().String("tech-postfix", "", "Tech postfix")
	cmd.Flags().String("type", "sip", "Provider type")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("hostname")
	return cmd
}

func newProvidersUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			hostname, _ := cmd.Flags().GetString("hostname")
			detail, _ := cmd.Flags().GetString("detail")
			techPrefix, _ := cmd.Flags().GetString("tech-prefix")
			techPostfix, _ := cmd.Flags().GetString("tech-postfix")
			provType, _ := cmd.Flags().GetString("type")

			body := voipbin_client.PutProvidersIdJSONRequestBody{
				Name:        name,
				Hostname:    hostname,
				Detail:      detail,
				TechPrefix:  techPrefix,
				TechPostfix: techPostfix,
				Type:        voipbin_client.RouteManagerProviderType(provType),
				TechHeaders: map[string]interface{}{},
			}

			resp, err := client.PutProvidersIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update provider: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, providerDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Provider name")
	cmd.Flags().String("hostname", "", "Destination hostname")
	cmd.Flags().String("detail", "", "Provider detail")
	cmd.Flags().String("tech-prefix", "", "Tech prefix")
	cmd.Flags().String("tech-postfix", "", "Tech postfix")
	cmd.Flags().String("type", "", "Provider type")
	return cmd
}

func newProvidersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteProvidersIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete provider: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Provider %s deleted.\n", args[0])
			return nil
		},
	}
}
