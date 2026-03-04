package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newConferencesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conferences",
		Short: "Manage conferences",
	}
	cmd.AddCommand(
		newConferencesListCmd(),
		newConferencesGetCmd(),
		newConferencesCreateCmd(),
		newConferencesUpdateCmd(),
		newConferencesDeleteCmd(),
		newConferencesRecordingStartCmd(),
		newConferencesRecordingStopCmd(),
		newConferencesTranscribeStartCmd(),
		newConferencesTranscribeStopCmd(),
	)
	return cmd
}

var conferenceListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "STATUS", Field: "status"},
	{Name: "TYPE", Field: "type"},
	{Name: "RECORDING_ID", Field: "recording_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var conferenceDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "STATUS", Field: "status"},
	{Name: "TYPE", Field: "type"},
	{Name: "RECORDING_ID", Field: "recording_id"},
	{Name: "TRANSCRIBE_ID", Field: "transcribe_id"},
	{Name: "PRE_FLOW_ID", Field: "pre_flow_id"},
	{Name: "POST_FLOW_ID", Field: "post_flow_id"},
	{Name: "TIMEOUT", Field: "timeout"},
	{Name: "CREATED", Field: "tm_create"},
}

func newConferencesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List conferences",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetConferencesParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetConferencesWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list conferences: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, conferenceListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newConferencesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a conference by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetConferencesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get conference: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, conferenceDetailColumns)
		},
	}
}

func newConferencesCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new conference",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			preFlowID, _ := cmd.Flags().GetString("pre-flow-id")
			postFlowID, _ := cmd.Flags().GetString("post-flow-id")
			timeout, _ := cmd.Flags().GetInt("timeout")
			confType, _ := cmd.Flags().GetString("type")

			body := voipbin_client.PostConferencesJSONRequestBody{
				Name:       name,
				Detail:     detail,
				PreFlowId:  preFlowID,
				PostFlowId: postFlowID,
				Timeout:    timeout,
				Type:       voipbin_client.ConferenceManagerConferenceType(confType),
				Data:       map[string]any{},
			}

			resp, err := client.PostConferencesWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create conference: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, conferenceDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Conference name")
	cmd.Flags().String("detail", "", "Conference detail")
	cmd.Flags().String("pre-flow-id", "", "Pre-flow ID")
	cmd.Flags().String("post-flow-id", "", "Post-flow ID")
	cmd.Flags().Int("timeout", 0, "Conference timeout in seconds")
	cmd.Flags().String("type", "", "Conference type")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newConferencesUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a conference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			preFlowID, _ := cmd.Flags().GetString("pre-flow-id")
			postFlowID, _ := cmd.Flags().GetString("post-flow-id")
			timeout, _ := cmd.Flags().GetInt("timeout")

			body := voipbin_client.PutConferencesIdJSONRequestBody{
				Name:       name,
				Detail:     detail,
				PreFlowId:  preFlowID,
				PostFlowId: postFlowID,
				Timeout:    timeout,
				Data:       map[string]any{},
			}

			resp, err := client.PutConferencesIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update conference: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, conferenceDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Conference name")
	cmd.Flags().String("detail", "", "Conference detail")
	cmd.Flags().String("pre-flow-id", "", "Pre-flow ID")
	cmd.Flags().String("post-flow-id", "", "Post-flow ID")
	cmd.Flags().Int("timeout", 0, "Conference timeout in seconds")
	return cmd
}

func newConferencesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a conference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteConferencesIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete conference: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conference %s deleted.\n", args[0])
			return nil
		},
	}
}

func newConferencesRecordingStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recording-start <id>",
		Short: "Start recording a conference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			duration, _ := cmd.Flags().GetInt("duration")
			format, _ := cmd.Flags().GetString("format")
			onEndFlowID, _ := cmd.Flags().GetString("on-end-flow-id")

			body := voipbin_client.PostConferencesIdRecordingStartJSONRequestBody{
				Duration:    duration,
				Format:      voipbin_client.PostConferencesIdRecordingStartJSONBodyFormat(format),
				OnEndFlowId: onEndFlowID,
			}

			resp, err := client.PostConferencesIdRecordingStartWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not start recording: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conference %s recording started.\n", args[0])
			return nil
		},
	}
	cmd.Flags().Int("duration", 0, "Maximum recording duration in seconds")
	cmd.Flags().String("format", "", "Recording format")
	cmd.Flags().String("on-end-flow-id", "", "Flow ID to execute when recording ends")
	return cmd
}

func newConferencesRecordingStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recording-stop <id>",
		Short: "Stop recording a conference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.PostConferencesIdRecordingStopWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not stop recording: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conference %s recording stopped.\n", args[0])
			return nil
		},
	}
}

func newConferencesTranscribeStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transcribe-start <id>",
		Short: "Start transcribing a conference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			language, _ := cmd.Flags().GetString("language")

			body := voipbin_client.PostConferencesIdTranscribeStartJSONRequestBody{
				Language: language,
			}

			resp, err := client.PostConferencesIdTranscribeStartWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not start transcription: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conference %s transcription started.\n", args[0])
			return nil
		},
	}
	cmd.Flags().String("language", "", "Transcription language")
	return cmd
}

func newConferencesTranscribeStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "transcribe-stop <id>",
		Short: "Stop transcribing a conference",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.PostConferencesIdTranscribeStopWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not stop transcription: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Conference %s transcription stopped.\n", args[0])
			return nil
		},
	}
}
