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

			items, nextToken, err := c.List(context.Background(), "/ais", params)
			if err != nil {
				return fmt.Errorf("could not list AIs: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, aiListColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/ais/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get AI: %w", err)
			}

			return output.PrintItem(cmd, result, aiDetailColumns)
		},
	}
}

func newAisCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
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

			body := map[string]interface{}{
				"name":         name,
				"detail":       detail,
				"engine_type":  engineType,
				"engine_model": engineModel,
				"engine_key":   engineKey,
				"engine_data":  map[string]interface{}{},
				"init_prompt":  initPrompt,
				"stt_type":     sttType,
				"tts_type":     ttsType,
				"tts_voice_id": ttsVoiceID,
			}

			result, err := c.Post(context.Background(), "/ais", body)
			if err != nil {
				return fmt.Errorf("could not create AI: %w", err)
			}

			return output.PrintItem(cmd, result, aiDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
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

			body := map[string]interface{}{
				"name":         name,
				"detail":       detail,
				"engine_type":  engineType,
				"engine_model": engineModel,
				"engine_key":   engineKey,
				"engine_data":  map[string]interface{}{},
				"init_prompt":  initPrompt,
				"stt_type":     sttType,
				"tts_type":     ttsType,
				"tts_voice_id": ttsVoiceID,
			}

			result, err := c.Put(context.Background(), "/ais/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update AI: %w", err)
			}

			return output.PrintItem(cmd, result, aiDetailColumns)
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/ais/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI %s deleted.\n", args[0])
			return nil
		},
	}
}
