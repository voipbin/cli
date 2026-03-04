package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newCampaignsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campaigns",
		Short: "Manage campaigns",
	}
	cmd.AddCommand(
		newCampaignsListCmd(),
		newCampaignsGetCmd(),
		newCampaignsCreateCmd(),
		newCampaignsUpdateCmd(),
		newCampaignsDeleteCmd(),
		newCampaignsSetStatusCmd(),
		newCampaignsSetServiceLevelCmd(),
		newCampaignsSetActionsCmd(),
		newCampaignsSetNextCampaignCmd(),
		newCampaignsSetResourceInfoCmd(),
	)
	return cmd
}

var campaignListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "STATUS", Field: "status"},
	{Name: "TYPE", Field: "type"},
	{Name: "OUTPLAN_ID", Field: "outplan_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var campaignDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "STATUS", Field: "status"},
	{Name: "TYPE", Field: "type"},
	{Name: "END_HANDLE", Field: "end_handle"},
	{Name: "OUTPLAN_ID", Field: "outplan_id"},
	{Name: "OUTDIAL_ID", Field: "outdial_id"},
	{Name: "QUEUE_ID", Field: "queue_id"},
	{Name: "SERVICE_LEVEL", Field: "service_level"},
	{Name: "NEXT_CAMPAIGN_ID", Field: "next_campaign_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newCampaignsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List campaigns",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetCampaignsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetCampaignsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list campaigns: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, campaignListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newCampaignsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a campaign by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetCampaignsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get campaign: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, campaignDetailColumns)
		},
	}
}

func newCampaignsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			campType, _ := cmd.Flags().GetString("type")
			endHandle, _ := cmd.Flags().GetString("end-handle")
			outplanID, _ := cmd.Flags().GetString("outplan-id")
			outdialID, _ := cmd.Flags().GetString("outdial-id")
			queueID, _ := cmd.Flags().GetString("queue-id")
			nextCampaignID, _ := cmd.Flags().GetString("next-campaign-id")
			serviceLevel, _ := cmd.Flags().GetInt("service-level")

			body := voipbin_client.PostCampaignsJSONRequestBody{
				Name:           name,
				Detail:         detail,
				Type:           voipbin_client.CampaignManagerCampaignType(campType),
				EndHandle:      voipbin_client.CampaignManagerCampaignEndHandle(endHandle),
				OutplanId:      outplanID,
				OutdialId:      outdialID,
				QueueId:        queueID,
				NextCampaignId: nextCampaignID,
				ServiceLevel:   serviceLevel,
				Actions:        []voipbin_client.FlowManagerAction{},
			}

			resp, err := client.PostCampaignsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create campaign: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, campaignDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Campaign name")
	cmd.Flags().String("detail", "", "Campaign detail")
	cmd.Flags().String("type", "", "Campaign type")
	cmd.Flags().String("end-handle", "", "End handle behavior")
	cmd.Flags().String("outplan-id", "", "Outplan ID")
	cmd.Flags().String("outdial-id", "", "Outdial ID")
	cmd.Flags().String("queue-id", "", "Queue ID")
	cmd.Flags().String("next-campaign-id", "", "Next campaign ID")
	cmd.Flags().Int("service-level", 0, "Service level")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}

func newCampaignsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a campaign",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			campType, _ := cmd.Flags().GetString("type")
			endHandle, _ := cmd.Flags().GetString("end-handle")
			serviceLevel, _ := cmd.Flags().GetInt("service-level")

			body := voipbin_client.PutCampaignsIdJSONRequestBody{
				Name:         name,
				Detail:       detail,
				Type:         voipbin_client.CampaignManagerCampaignType(campType),
				EndHandle:    voipbin_client.CampaignManagerCampaignEndHandle(endHandle),
				ServiceLevel: serviceLevel,
			}

			resp, err := client.PutCampaignsIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update campaign: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, campaignDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Campaign name")
	cmd.Flags().String("detail", "", "Campaign detail")
	cmd.Flags().String("type", "", "Campaign type")
	cmd.Flags().String("end-handle", "", "End handle behavior")
	cmd.Flags().Int("service-level", 0, "Service level")
	return cmd
}

func newCampaignsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a campaign",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteCampaignsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete campaign: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s deleted.\n", args[0])
			return nil
		},
	}
}

func newCampaignsSetStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-status <id>",
		Short: "Set campaign status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			body := voipbin_client.PutCampaignsIdStatusJSONRequestBody{
				Status: voipbin_client.CampaignManagerCampaignStatus(status),
			}

			resp, err := client.PutCampaignsIdStatusWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set campaign status: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s status set to %s.\n", args[0], status)
			return nil
		},
	}
	cmd.Flags().String("status", "", "Campaign status")
	_ = cmd.MarkFlagRequired("status")
	return cmd
}

func newCampaignsSetServiceLevelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-service-level <id>",
		Short: "Set campaign service level",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			serviceLevel, _ := cmd.Flags().GetInt("service-level")

			body := voipbin_client.PutCampaignsIdServiceLevelJSONRequestBody{
				ServiceLevel: serviceLevel,
			}

			resp, err := client.PutCampaignsIdServiceLevelWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set service level: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s service level updated.\n", args[0])
			return nil
		},
	}
	cmd.Flags().Int("service-level", 0, "Service level")
	_ = cmd.MarkFlagRequired("service-level")
	return cmd
}

func newCampaignsSetActionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-actions <id>",
		Short: "Set campaign actions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			body := voipbin_client.PutCampaignsIdActionsJSONRequestBody{
				Actions: []voipbin_client.FlowManagerAction{},
			}

			resp, err := client.PutCampaignsIdActionsWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set campaign actions: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s actions updated.\n", args[0])
			return nil
		},
	}
	return cmd
}

func newCampaignsSetNextCampaignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-next-campaign <id>",
		Short: "Set next campaign ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			nextCampaignID, _ := cmd.Flags().GetString("next-campaign-id")

			body := voipbin_client.PutCampaignsIdNextCampaignIdJSONRequestBody{
				NextCampaignId: nextCampaignID,
			}

			resp, err := client.PutCampaignsIdNextCampaignIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set next campaign: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s next campaign set.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("next-campaign-id", "", "Next campaign ID")
	_ = cmd.MarkFlagRequired("next-campaign-id")
	return cmd
}

func newCampaignsSetResourceInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-resource-info <id>",
		Short: "Set campaign resource info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			outplanID, _ := cmd.Flags().GetString("outplan-id")
			outdialID, _ := cmd.Flags().GetString("outdial-id")
			queueID, _ := cmd.Flags().GetString("queue-id")
			nextCampaignID, _ := cmd.Flags().GetString("next-campaign-id")

			body := voipbin_client.PutCampaignsIdResourceInfoJSONRequestBody{
				OutplanId:      outplanID,
				OutdialId:      outdialID,
				QueueId:        queueID,
				NextCampaignId: nextCampaignID,
			}

			resp, err := client.PutCampaignsIdResourceInfoWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not set resource info: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s resource info updated.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("outplan-id", "", "Outplan ID")
	cmd.Flags().String("outdial-id", "", "Outdial ID")
	cmd.Flags().String("queue-id", "", "Queue ID")
	cmd.Flags().String("next-campaign-id", "", "Next campaign ID")
	return cmd
}
