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

func newTranscriptsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transcripts",
		Short: "View transcripts",
	}
	cmd.AddCommand(
		newTranscriptsListCmd(),
	)
	return cmd
}

var transcriptListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "TRANSCRIBE_ID", Field: "transcribe_id"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "MESSAGE", Field: "message"},
	{Name: "TM_TRANSCRIPT", Field: "tm_transcript"},
	{Name: "CREATED", Field: "tm_create"},
}

func newTranscriptsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List transcripts for a transcribe",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			transcribeID, _ := cmd.Flags().GetString("transcribe-id")
			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := url.Values{}
			params.Set("transcribe_id", transcribeID)
			if pageToken != "" {
				params.Set("page_token", pageToken)
			}
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}

			items, nextToken, err := c.List(context.Background(), "/transcripts", params)
			if err != nil {
				return fmt.Errorf("could not list transcripts: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, transcriptListColumns)
		},
	}
	cmd.Flags().String("transcribe-id", "", "Transcribe ID to list transcripts for")
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	_ = cmd.MarkFlagRequired("transcribe-id")
	return cmd
}
