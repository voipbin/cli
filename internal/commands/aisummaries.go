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

func newAisummariesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aisummaries",
		Short: "Manage AI summaries",
	}
	cmd.AddCommand(
		newAisummariesListCmd(),
		newAisummariesGetCmd(),
		newAisummariesCreateCmd(),
		newAisummariesDeleteCmd(),
	)
	return cmd
}

var aisummaryListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "CREATED", Field: "tm_create"},
}

var aisummaryDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "CONTENT", Field: "content"},
	{Name: "ACTIVEFLOW_ID", Field: "activeflow_id"},
	{Name: "ON_END_FLOW_ID", Field: "on_end_flow_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newAisummariesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List AI summaries",
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

			items, nextToken, err := c.List(context.Background(), "/aisummaries", params)
			if err != nil {
				return fmt.Errorf("could not list AI summaries: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, aisummaryListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newAisummariesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an AI summary by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/aisummaries/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get AI summary: %w", err)
			}

			return output.PrintItem(cmd, result, aisummaryDetailColumns)
		},
	}
}

func newAisummariesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			referenceID, _ := cmd.Flags().GetString("reference-id")
			referenceType, _ := cmd.Flags().GetString("reference-type")
			language, _ := cmd.Flags().GetString("language")
			onEndFlowID, _ := cmd.Flags().GetString("on-end-flow-id")

			body := map[string]interface{}{
				"reference_id":    referenceID,
				"reference_type":  referenceType,
				"language":        language,
				"on_end_flow_id":  onEndFlowID,
			}

			result, err := c.Post(context.Background(), "/aisummaries", body)
			if err != nil {
				return fmt.Errorf("could not create AI summary: %w", err)
			}

			return output.PrintItem(cmd, result, aisummaryDetailColumns)
		},
	}
	cmd.Flags().String("reference-id", "", "Reference ID")
	cmd.Flags().String("reference-type", "", "Reference type")
	cmd.Flags().String("language", "", "Summary language")
	cmd.Flags().String("on-end-flow-id", "", "Flow ID to execute when summary ends")
	_ = cmd.MarkFlagRequired("reference-id")
	_ = cmd.MarkFlagRequired("reference-type")
	return cmd
}

func newAisummariesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an AI summary",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/aisummaries/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI summary: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI summary %s deleted.\n", args[0])
			return nil
		},
	}
}
