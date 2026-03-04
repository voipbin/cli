package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newOutdialsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outdials",
		Short: "Manage outdials",
	}
	cmd.AddCommand(
		newOutdialsListCmd(),
		newOutdialsGetCmd(),
		newOutdialsCreateCmd(),
		newOutdialsUpdateCmd(),
		newOutdialsDeleteCmd(),
		newOutdialsSetCampaignCmd(),
		newOutdialsSetDataCmd(),
		newOutdialsListTargetsCmd(),
		newOutdialsCreateTargetCmd(),
		newOutdialsDeleteTargetCmd(),
	)
	return cmd
}

var outdialListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "CAMPAIGN_ID", Field: "campaign_id"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CREATED", Field: "tm_create"},
}

var outdialDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "CAMPAIGN_ID", Field: "campaign_id"},
	{Name: "DATA", Field: "data"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newOutdialsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List outdials",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetOutdialsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetOutdialsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list outdials: %w", err)
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
			return output.PrintList(cmd, *resp.JSON200.Result, outdialListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newOutdialsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an outdial by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.GetOutdialsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get outdial: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outdialDetailColumns)
		},
	}
}

func newOutdialsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new outdial",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			campaignID, _ := cmd.Flags().GetString("campaign-id")
			data, _ := cmd.Flags().GetString("data")
			body := voipbin_client.PostOutdialsJSONRequestBody{
				Name:       name,
				Detail:     detail,
				CampaignId: campaignID,
				Data:       data,
			}
			resp, err := client.PostOutdialsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create outdial: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outdialDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Outdial name")
	cmd.Flags().String("detail", "", "Outdial detail")
	cmd.Flags().String("campaign-id", "", "Campaign ID")
	cmd.Flags().String("data", "", "Outdial data")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("campaign-id")
	return cmd
}

func newOutdialsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an outdial",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			body := voipbin_client.PutOutdialsIdJSONRequestBody{
				Name:   name,
				Detail: detail,
			}
			resp, err := client.PutOutdialsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update outdial: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}
			return output.PrintItem(cmd, resp.JSON200, outdialDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "New name")
	cmd.Flags().String("detail", "", "New detail")
	return cmd
}

func newOutdialsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an outdial",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteOutdialsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete outdial: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Outdial %s deleted.\n", args[0])
			return nil
		},
	}
}

func newOutdialsSetCampaignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-campaign <id>",
		Short: "Set campaign ID for an outdial",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			campaignID, _ := cmd.Flags().GetString("campaign-id")
			body := voipbin_client.PutOutdialsIdCampaignIdJSONRequestBody{CampaignId: campaignID}
			resp, err := client.PutOutdialsIdCampaignIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set campaign: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Outdial %s campaign set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("campaign-id", "", "Campaign ID")
	_ = cmd.MarkFlagRequired("campaign-id")
	return cmd
}

func newOutdialsSetDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-data <id>",
		Short: "Set data for an outdial",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			data, _ := cmd.Flags().GetString("data")
			body := voipbin_client.PutOutdialsIdDataJSONRequestBody{Data: data}
			resp, err := client.PutOutdialsIdDataWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set data: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Outdial %s data set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("data", "", "Data value")
	_ = cmd.MarkFlagRequired("data")
	return cmd
}

func newOutdialsListTargetsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-targets <id>",
		Short: "List targets for an outdial",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")
			params := &voipbin_client.GetOutdialsIdTargetsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}
			resp, err := client.GetOutdialsIdTargetsWithResponse(context.Background(), args[0], params)
			if err != nil {
				return fmt.Errorf("could not list targets: %w", err)
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
			cols := []output.Column{
				{Name: "ID", Field: "id"},
				{Name: "NAME", Field: "name"},
				{Name: "DETAIL", Field: "detail"},
				{Name: "DATA", Field: "data"},
			}
			return output.PrintList(cmd, *resp.JSON200.Result, cols)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newOutdialsCreateTargetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-target <id>",
		Short: "Create a target for an outdial",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			data, _ := cmd.Flags().GetString("data")
			dest0, _ := cmd.Flags().GetString("destination-0")
			body := voipbin_client.PostOutdialsIdTargetsJSONRequestBody{
				Name:         name,
				Detail:       detail,
				Data:         data,
				Destination0: voipbin_client.CommonAddress{Target: &dest0},
			}
			resp, err := client.PostOutdialsIdTargetsWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not create target: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Target created for outdial %s.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("name", "", "Target name")
	cmd.Flags().String("detail", "", "Target detail")
	cmd.Flags().String("data", "", "Target data")
	cmd.Flags().String("destination-0", "", "Primary destination target")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("destination-0")
	return cmd
}

func newOutdialsDeleteTargetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-target <id> <target-id>",
		Short: "Delete a target from an outdial",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}
			resp, err := client.DeleteOutdialsIdTargetsTargetIdWithResponse(context.Background(), args[0], args[1])
			if err != nil {
				return fmt.Errorf("could not delete target: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Target %s deleted from outdial %s.\n", args[1], args[0])
			return nil
		},
	}
}
