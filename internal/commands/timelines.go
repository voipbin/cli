package commands

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
)

func newTimelinesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "timelines",
		Short: "View timeline events and analysis",
	}
	cmd.AddCommand(
		newTimelinesListEventsCmd(),
		newTimelinesAggregatedEventsCmd(),
		newTimelinesSipAnalysisCmd(),
		newTimelinesPcapCmd(),
	)
	return cmd
}

var timelineEventColumns = []output.Column{
	{Name: "TYPE", Field: "type"},
	{Name: "TIMESTAMP", Field: "timestamp"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newTimelinesListEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-events",
		Short: "List timeline events for a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resourceType, _ := cmd.Flags().GetString("resource-type")
			resourceID, _ := cmd.Flags().GetString("resource-id")
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := url.Values{}
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}

			items, nextToken, err := c.List(context.Background(), "/timelines/"+resourceType+"/"+resourceID, params)
			if err != nil {
				return fmt.Errorf("could not list timeline events: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, timelineEventColumns)
		},
	}
	cmd.Flags().String("resource-type", "", "Resource type (e.g., call, activeflow)")
	cmd.Flags().String("resource-id", "", "Resource ID")
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	_ = cmd.MarkFlagRequired("resource-type")
	_ = cmd.MarkFlagRequired("resource-id")
	return cmd
}

func newTimelinesAggregatedEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregated-events",
		Short: "List aggregated timeline events",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			activeflowID, _ := cmd.Flags().GetString("activeflow-id")
			callID, _ := cmd.Flags().GetString("call-id")
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := url.Values{}
			if activeflowID != "" {
				params.Set("activeflow_id", activeflowID)
			}
			if callID != "" {
				params.Set("call_id", callID)
			}
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}

			items, nextToken, err := c.List(context.Background(), "/timelines/aggregated_events", params)
			if err != nil {
				return fmt.Errorf("could not list aggregated events: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, timelineEventColumns)
		},
	}
	cmd.Flags().String("activeflow-id", "", "Activeflow ID to filter by")
	cmd.Flags().String("call-id", "", "Call ID to filter by")
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newTimelinesSipAnalysisCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sip-analysis <call-id>",
		Short: "Get SIP analysis for a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/timelines/call_sip_analysis/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get SIP analysis: %w", err)
			}

			columns := []output.Column{
				{Name: "CALL_ID", Field: "call_id"},
				{Name: "SIP_MESSAGES", Field: "sip_messages"},
			}

			return output.PrintItem(cmd, result, columns)
		},
	}
}

func newTimelinesPcapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pcap <call-id>",
		Short: "Download PCAP file for a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := c.RawGet(context.Background(), "/timelines/call_pcap/"+args[0])
			if err != nil {
				return fmt.Errorf("could not download PCAP: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
			}

			outputFile, _ := cmd.Flags().GetString("output-file")
			var w io.Writer
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("could not create file: %w", err)
				}
				defer f.Close()
				w = f
			} else {
				w = cmd.OutOrStdout()
			}

			n, err := io.Copy(w, resp.Body)
			if err != nil {
				return fmt.Errorf("could not write PCAP: %w", err)
			}

			if outputFile != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "PCAP saved to %s (%d bytes)\n", outputFile, n)
			}

			return nil
		},
	}
	cmd.Flags().String("output-file", "", "Output file path (default: stdout)")
	return cmd
}
