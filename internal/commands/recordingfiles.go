package commands

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
)

func newRecordingfilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recordingfiles",
		Short: "Access recording files",
	}
	cmd.AddCommand(
		newRecordingfilesGetCmd(),
	)
	return cmd
}

func newRecordingfilesGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get a recording file (follows redirect to download URL)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := c.RawGet(context.Background(), "/recordingfiles/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get recording file: %w", err)
			}
			defer resp.Body.Close()

			// For 307 redirects, the CheckRedirect policy blocks cross-domain redirects,
			// so we get the last response with the Location header.
			if resp.StatusCode == 307 || resp.StatusCode == 302 || resp.StatusCode == 301 {
				location := resp.Header.Get("Location")
				if location != "" {
					fmt.Fprintln(cmd.OutOrStdout(), location)
					return nil
				}
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				body, _ := io.ReadAll(resp.Body)
				return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
			}

			// If we got a direct response (e.g., same-domain), stream the file
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
				return fmt.Errorf("could not write recording: %w", err)
			}

			if outputFile != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Recording saved to %s (%d bytes)\n", outputFile, n)
			}

			return nil
		},
	}
	cmd.Flags().String("output-file", "", "Output file path (default: stdout)")
	return cmd
}
