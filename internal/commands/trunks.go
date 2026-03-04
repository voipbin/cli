package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newTrunksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trunks",
		Short: "Manage trunks",
	}
	cmd.AddCommand(
		newTrunksListCmd(),
		newTrunksGetCmd(),
		newTrunksCreateCmd(),
		newTrunksUpdateCmd(),
		newTrunksDeleteCmd(),
	)
	return cmd
}

var trunkListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DOMAIN_NAME", Field: "domain_name"},
	{Name: "USERNAME", Field: "username"},
	{Name: "CREATED", Field: "tm_create"},
}

var trunkDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "DOMAIN_NAME", Field: "domain_name"},
	{Name: "USERNAME", Field: "username"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newTrunksListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List trunks",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetTrunksParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetTrunksWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list trunks: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, trunkListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newTrunksGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a trunk by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetTrunksIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get trunk: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, trunkDetailColumns)
		},
	}
}

func newTrunksCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new trunk",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			domainName, _ := cmd.Flags().GetString("domain-name")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PostTrunksJSONRequestBody{
				Name:       name,
				DomainName: domainName,
				Username:   username,
				Password:   password,
				Detail:     detail,
				AllowedIps: []string{},
				AuthTypes:  []voipbin_client.RegistrarManagerAuthType{},
			}

			resp, err := client.PostTrunksWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create trunk: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, trunkDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Trunk name")
	cmd.Flags().String("domain-name", "", "Domain name")
	cmd.Flags().String("username", "", "Username")
	cmd.Flags().String("password", "", "Password")
	cmd.Flags().String("detail", "", "Trunk detail")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newTrunksUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a trunk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			detail, _ := cmd.Flags().GetString("detail")

			body := voipbin_client.PutTrunksIdJSONRequestBody{
				Name:       name,
				Username:   username,
				Password:   password,
				Detail:     detail,
				AllowedIps: []string{},
				AuthTypes:  []voipbin_client.RegistrarManagerAuthType{},
			}

			resp, err := client.PutTrunksIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update trunk: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, trunkDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Trunk name")
	cmd.Flags().String("username", "", "Username")
	cmd.Flags().String("password", "", "Password")
	cmd.Flags().String("detail", "", "Trunk detail")
	return cmd
}

func newTrunksDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a trunk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteTrunksIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete trunk: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Trunk %s deleted.\n", args[0])
			return nil
		},
	}
}
