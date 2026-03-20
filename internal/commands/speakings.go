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

func newSpeakingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "speakings",
		Short: "Manage text-to-speech speakings",
	}
	cmd.AddCommand(
		newSpeakingsListCmd(),
		newSpeakingsGetCmd(),
		newSpeakingsCreateCmd(),
		newSpeakingsDeleteCmd(),
		newSpeakingsSayCmd(),
		newSpeakingsFlushCmd(),
		newSpeakingsStopCmd(),
	)
	return cmd
}

var speakingListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "STATUS", Field: "status"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "PROVIDER", Field: "provider"},
	{Name: "CREATED", Field: "tm_create"},
}

var speakingDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "PROVIDER", Field: "provider"},
	{Name: "VOICE_ID", Field: "voice_id"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newSpeakingsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List speakings",
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

			items, nextToken, err := c.List(context.Background(), "/speakings", params)
			if err != nil {
				return fmt.Errorf("could not list speakings: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, speakingListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newSpeakingsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a speaking by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/speakings/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get speaking: %w", err)
			}

			return output.PrintItem(cmd, result, speakingDetailColumns)
		},
	}
}

func newSpeakingsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new speaking",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			referenceType, _ := cmd.Flags().GetString("reference-type")
			referenceID, _ := cmd.Flags().GetString("reference-id")
			language, _ := cmd.Flags().GetString("language")
			provider, _ := cmd.Flags().GetString("provider")
			voiceID, _ := cmd.Flags().GetString("voice-id")
			direction, _ := cmd.Flags().GetString("direction")

			body := map[string]interface{}{
				"reference_type": referenceType,
				"reference_id":   referenceID,
			}
			if language != "" {
				body["language"] = language
			}
			if provider != "" {
				body["provider"] = provider
			}
			if voiceID != "" {
				body["voice_id"] = voiceID
			}
			if direction != "" {
				body["direction"] = direction
			}

			result, err := c.Post(context.Background(), "/speakings", body)
			if err != nil {
				return fmt.Errorf("could not create speaking: %w", err)
			}

			return output.PrintItem(cmd, result, speakingDetailColumns)
		},
	}
	cmd.Flags().String("reference-type", "", "Reference type")
	cmd.Flags().String("reference-id", "", "Reference ID")
	cmd.Flags().String("language", "", "Language code")
	cmd.Flags().String("provider", "", "TTS provider")
	cmd.Flags().String("voice-id", "", "Voice ID")
	cmd.Flags().String("direction", "", "Direction")
	_ = cmd.MarkFlagRequired("reference-type")
	_ = cmd.MarkFlagRequired("reference-id")
	return cmd
}

func newSpeakingsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a speaking",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/speakings/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete speaking: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Speaking %s deleted.\n", args[0])
			return nil
		},
	}
}

func newSpeakingsSayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "say <id>",
		Short: "Say text on a speaking",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			text, _ := cmd.Flags().GetString("text")
			body := map[string]interface{}{
				"text": text,
			}

			result, err := c.Post(context.Background(), "/speakings/"+args[0]+"/say", body)
			if err != nil {
				return fmt.Errorf("could not say text: %w", err)
			}

			return output.PrintItem(cmd, result, speakingDetailColumns)
		},
	}
	cmd.Flags().String("text", "", "Text to speak")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}

func newSpeakingsFlushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "flush <id>",
		Short: "Flush speaking queue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Post(context.Background(), "/speakings/"+args[0]+"/flush", nil)
			if err != nil {
				return fmt.Errorf("could not flush speaking: %w", err)
			}

			return output.PrintItem(cmd, result, speakingDetailColumns)
		},
	}
}

func newSpeakingsStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop <id>",
		Short: "Stop a speaking",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Post(context.Background(), "/speakings/"+args[0]+"/stop", nil)
			if err != nil {
				return fmt.Errorf("could not stop speaking: %w", err)
			}

			return output.PrintItem(cmd, result, speakingDetailColumns)
		},
	}
}
