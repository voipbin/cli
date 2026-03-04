package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetAisummariesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetAisummariesWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list AI summaries: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, aisummaryListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetAisummariesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get AI summary: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aisummaryDetailColumns)
		},
	}
}

func newAisummariesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			referenceID, _ := cmd.Flags().GetString("reference-id")
			referenceType, _ := cmd.Flags().GetString("reference-type")
			language, _ := cmd.Flags().GetString("language")
			onEndFlowID, _ := cmd.Flags().GetString("on-end-flow-id")

			body := voipbin_client.PostAisummariesJSONRequestBody{
				ReferenceId:   referenceID,
				ReferenceType: voipbin_client.AIManagerSummaryReferenceType(referenceType),
				Language:      language,
				OnEndFlowId:   onEndFlowID,
			}

			resp, err := client.PostAisummariesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create AI summary: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aisummaryDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteAisummariesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI summary: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI summary %s deleted.\n", args[0])
			return nil
		},
	}
}
