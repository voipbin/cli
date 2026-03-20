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

func newTranscribesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transcribes",
		Short: "Manage transcriptions",
	}
	cmd.AddCommand(
		newTranscribesListCmd(),
		newTranscribesGetCmd(),
		newTranscribesCreateCmd(),
		newTranscribesDeleteCmd(),
		newTranscribesStopCmd(),
	)
	return cmd
}

var transcribeListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "CREATED", Field: "tm_create"},
}

var transcribeDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "STATUS", Field: "status"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "CREATED", Field: "tm_create"},
}

func newTranscribesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List transcriptions",
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

			items, nextToken, err := c.List(context.Background(), "/transcribes", params)
			if err != nil {
				return fmt.Errorf("could not list transcribes: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, transcribeListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newTranscribesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a transcription by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/transcribes/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get transcribe: %w", err)
			}

			return output.PrintItem(cmd, result, transcribeDetailColumns)
		},
	}
}

func newTranscribesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new transcription",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			referenceID, _ := cmd.Flags().GetString("reference-id")
			referenceType, _ := cmd.Flags().GetString("reference-type")
			language, _ := cmd.Flags().GetString("language")
			direction, _ := cmd.Flags().GetString("direction")
			onEndFlowID, _ := cmd.Flags().GetString("on-end-flow-id")

			body := map[string]interface{}{
				"reference_id":    referenceID,
				"reference_type":  referenceType,
				"language":        language,
				"direction":       direction,
				"on_end_flow_id":  onEndFlowID,
			}

			result, err := c.Post(context.Background(), "/transcribes", body)
			if err != nil {
				return fmt.Errorf("could not create transcribe: %w", err)
			}

			return output.PrintItem(cmd, result, transcribeDetailColumns)
		},
	}
	cmd.Flags().String("reference-id", "", "Reference ID (call/conference/recording ID)")
	cmd.Flags().String("reference-type", "", "Reference type")
	cmd.Flags().String("language", "en-US", "BCP47 language code")
	cmd.Flags().String("direction", "both", "Transcription direction")
	cmd.Flags().String("on-end-flow-id", "", "Flow ID to execute on transcription end")
	_ = cmd.MarkFlagRequired("reference-id")
	_ = cmd.MarkFlagRequired("reference-type")
	return cmd
}

func newTranscribesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a transcription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Delete(context.Background(), "/transcribes/"+args[0]); err != nil {
				return fmt.Errorf("could not delete transcribe: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Transcription %s deleted.\n", args[0])
			return nil
		},
	}
}

func newTranscribesStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <id>",
		Short: "Stop a transcription",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			if _, err := c.Post(context.Background(), "/transcribes/"+args[0]+"/stop", nil); err != nil {
				return fmt.Errorf("could not stop transcribe: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Transcription %s stopped.\n", args[0])
			return nil
		},
	}
}
