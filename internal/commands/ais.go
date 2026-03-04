package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newAisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ais",
		Short: "Manage AIs",
	}
	cmd.AddCommand(
		newAisListCmd(),
		newAisGetCmd(),
		newAisCreateCmd(),
		newAisUpdateCmd(),
		newAisDeleteCmd(),
	)
	return cmd
}

var aiListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "ENGINE_TYPE", Field: "engine_type"},
	{Name: "ENGINE_MODEL", Field: "engine_model"},
	{Name: "STT_TYPE", Field: "stt_type"},
	{Name: "CREATED", Field: "tm_create"},
}

var aiDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "ENGINE_TYPE", Field: "engine_type"},
	{Name: "ENGINE_MODEL", Field: "engine_model"},
	{Name: "STT_TYPE", Field: "stt_type"},
	{Name: "TTS_TYPE", Field: "tts_type"},
	{Name: "TTS_VOICE_ID", Field: "tts_voice_id"},
	{Name: "INIT_PROMPT", Field: "init_prompt"},
	{Name: "CREATED", Field: "tm_create"},
}

func newAisListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List AIs",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetAisParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetAisWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list AIs: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, aiListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newAisGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an AI by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetAisIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get AI: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aiDetailColumns)
		},
	}
}

func newAisCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			engineType, _ := cmd.Flags().GetString("engine-type")
			engineModel, _ := cmd.Flags().GetString("engine-model")
			engineKey, _ := cmd.Flags().GetString("engine-key")
			initPrompt, _ := cmd.Flags().GetString("init-prompt")
			sttType, _ := cmd.Flags().GetString("stt-type")
			ttsType, _ := cmd.Flags().GetString("tts-type")
			ttsVoiceID, _ := cmd.Flags().GetString("tts-voice-id")

			body := voipbin_client.PostAisJSONRequestBody{
				Name:        name,
				Detail:      detail,
				EngineType:  voipbin_client.AIManagerAIEngineType(engineType),
				EngineModel: voipbin_client.AIManagerAIEngineModel(engineModel),
				EngineKey:   engineKey,
				EngineData:  map[string]any{},
				InitPrompt:  initPrompt,
				SttType:     sttType,
				TtsType:     ttsType,
				TtsVoiceId:  ttsVoiceID,
			}

			resp, err := client.PostAisWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create AI: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aiDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "AI name")
	cmd.Flags().String("detail", "", "AI detail")
	cmd.Flags().String("engine-type", "", "Engine type")
	cmd.Flags().String("engine-model", "", "Engine model")
	cmd.Flags().String("engine-key", "", "Engine API key")
	cmd.Flags().String("init-prompt", "", "Initial prompt")
	cmd.Flags().String("stt-type", "", "Speech-to-text type")
	cmd.Flags().String("tts-type", "", "Text-to-speech type")
	cmd.Flags().String("tts-voice-id", "", "TTS voice ID")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("engine-type")
	_ = cmd.MarkFlagRequired("engine-model")
	return cmd
}

func newAisUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an AI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			engineType, _ := cmd.Flags().GetString("engine-type")
			engineModel, _ := cmd.Flags().GetString("engine-model")
			engineKey, _ := cmd.Flags().GetString("engine-key")
			initPrompt, _ := cmd.Flags().GetString("init-prompt")
			sttType, _ := cmd.Flags().GetString("stt-type")
			ttsType, _ := cmd.Flags().GetString("tts-type")
			ttsVoiceID, _ := cmd.Flags().GetString("tts-voice-id")

			body := voipbin_client.PutAisIdJSONRequestBody{
				Name:        name,
				Detail:      detail,
				EngineType:  voipbin_client.AIManagerAIEngineType(engineType),
				EngineModel: voipbin_client.AIManagerAIEngineModel(engineModel),
				EngineKey:   engineKey,
				EngineData:  map[string]any{},
				InitPrompt:  initPrompt,
				SttType:     sttType,
				TtsType:     ttsType,
				TtsVoiceId:  ttsVoiceID,
			}

			resp, err := client.PutAisIdWithResponse(context.Background(), args[0], body)
			if err != nil {
				return fmt.Errorf("could not update AI: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aiDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "AI name")
	cmd.Flags().String("detail", "", "AI detail")
	cmd.Flags().String("engine-type", "", "Engine type")
	cmd.Flags().String("engine-model", "", "Engine model")
	cmd.Flags().String("engine-key", "", "Engine API key")
	cmd.Flags().String("init-prompt", "", "Initial prompt")
	cmd.Flags().String("stt-type", "", "Speech-to-text type")
	cmd.Flags().String("tts-type", "", "Text-to-speech type")
	cmd.Flags().String("tts-voice-id", "", "TTS voice ID")
	return cmd
}

func newAisDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an AI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteAisIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI %s deleted.\n", args[0])
			return nil
		},
	}
}
