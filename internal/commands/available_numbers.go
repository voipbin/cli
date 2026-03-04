package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newAvailableNumbersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "available-numbers",
		Short: "List available numbers",
	}
	cmd.AddCommand(newAvailableNumbersListCmd())
	return cmd
}

func newAvailableNumbersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available phone numbers",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			countryCode, _ := cmd.Flags().GetString("country-code")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetAvailableNumbersParams{
				CountryCode: countryCode,
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetAvailableNumbersWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list available numbers: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil || resp.JSON200.Result == nil {
				return fmt.Errorf("unexpected empty response")
			}

			for _, n := range *resp.JSON200.Result {
				fmt.Fprintln(cmd.OutOrStdout(), string(n))
			}
			return nil
		},
	}
	cmd.Flags().String("country-code", "", "ISO country code (e.g. US)")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	_ = cmd.MarkFlagRequired("country-code")
	return cmd
}
