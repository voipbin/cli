package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetTranscribesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetTranscribesWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list transcribes: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, transcribeListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetTranscribesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get transcribe: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, transcribeDetailColumns)
		},
	}
}

func newTranscribesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new transcription",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			referenceID, _ := cmd.Flags().GetString("reference-id")
			referenceType, _ := cmd.Flags().GetString("reference-type")
			language, _ := cmd.Flags().GetString("language")
			direction, _ := cmd.Flags().GetString("direction")
			onEndFlowID, _ := cmd.Flags().GetString("on-end-flow-id")

			body := voipbin_client.PostTranscribesJSONRequestBody{
				ReferenceId:   referenceID,
				ReferenceType: voipbin_client.TranscribeManagerTranscribeReferenceType(referenceType),
				Language:      language,
				Direction:     voipbin_client.TranscribeManagerTranscribeDirection(direction),
				OnEndFlowId:   onEndFlowID,
			}

			resp, err := client.PostTranscribesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create transcribe: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, transcribeDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteTranscribesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete transcribe: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.PostTranscribesIdStopWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not stop transcribe: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Transcription %s stopped.\n", args[0])
			return nil
		},
	}
}
