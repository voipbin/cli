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

			items, nextToken, err := c.List(context.Background(), "/campaigns", params)
			if err != nil {
				return fmt.Errorf("could not list campaigns: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, campaignListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/campaigns/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get campaign: %w", err)
			}

			return output.PrintItem(cmd, result, campaignDetailColumns)
		},
	}
}

func newCampaignsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
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

			body := map[string]interface{}{
				"name":             name,
				"detail":           detail,
				"type":             campType,
				"end_handle":       endHandle,
				"outplan_id":       outplanID,
				"outdial_id":       outdialID,
				"queue_id":         queueID,
				"next_campaign_id": nextCampaignID,
				"service_level":    serviceLevel,
				"actions":          []interface{}{},
			}

			result, err := c.Post(context.Background(), "/campaigns", body)
			if err != nil {
				return fmt.Errorf("could not create campaign: %w", err)
			}

			return output.PrintItem(cmd, result, campaignDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			campType, _ := cmd.Flags().GetString("type")
			endHandle, _ := cmd.Flags().GetString("end-handle")
			serviceLevel, _ := cmd.Flags().GetInt("service-level")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if campType != "" {
				body["type"] = campType
			}
			if endHandle != "" {
				body["end_handle"] = endHandle
			}
			if serviceLevel > 0 {
				body["service_level"] = serviceLevel
			}

			result, err := c.Put(context.Background(), "/campaigns/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update campaign: %w", err)
			}

			return output.PrintItem(cmd, result, campaignDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/campaigns/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete campaign: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			status, _ := cmd.Flags().GetString("status")

			body := map[string]interface{}{
				"status": status,
			}

			_, err = c.Put(context.Background(), "/campaigns/"+args[0]+"/status", body)
			if err != nil {
				return fmt.Errorf("could not set campaign status: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			serviceLevel, _ := cmd.Flags().GetInt("service-level")

			body := map[string]interface{}{
				"service_level": serviceLevel,
			}

			_, err = c.Put(context.Background(), "/campaigns/"+args[0]+"/service_level", body)
			if err != nil {
				return fmt.Errorf("could not set service level: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			actionsJSON, _ := cmd.Flags().GetString("actions")
			var actions interface{}
			if err := json.Unmarshal([]byte(actionsJSON), &actions); err != nil {
				return fmt.Errorf("invalid actions JSON: %w", err)
			}

			body := map[string]interface{}{
				"actions": actions,
			}

			_, err = c.Put(context.Background(), "/campaigns/"+args[0]+"/actions", body)
			if err != nil {
				return fmt.Errorf("could not set campaign actions: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Campaign %s actions updated.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("actions", "", "Actions as JSON array")
	_ = cmd.MarkFlagRequired("actions")
	return cmd
}

func newCampaignsSetNextCampaignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-next-campaign <id>",
		Short: "Set next campaign ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			nextCampaignID, _ := cmd.Flags().GetString("next-campaign-id")

			body := map[string]interface{}{
				"next_campaign_id": nextCampaignID,
			}

			_, err = c.Put(context.Background(), "/campaigns/"+args[0]+"/next_campaign_id", body)
			if err != nil {
				return fmt.Errorf("could not set next campaign: %w", err)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			outplanID, _ := cmd.Flags().GetString("outplan-id")
			outdialID, _ := cmd.Flags().GetString("outdial-id")
			queueID, _ := cmd.Flags().GetString("queue-id")
			nextCampaignID, _ := cmd.Flags().GetString("next-campaign-id")

			body := map[string]interface{}{
				"outplan_id":       outplanID,
				"outdial_id":       outdialID,
				"queue_id":         queueID,
				"next_campaign_id": nextCampaignID,
			}

			_, err = c.Put(context.Background(), "/campaigns/"+args[0]+"/resource_info", body)
			if err != nil {
				return fmt.Errorf("could not set resource info: %w", err)
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
