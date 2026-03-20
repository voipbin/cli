package commands

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
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
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			countryCode, _ := cmd.Flags().GetString("country-code")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := url.Values{}
			params.Set("country_code", countryCode)
			if pageSize > 0 {
				params.Set("page_size", strconv.Itoa(pageSize))
			}

			items, _, err := c.List(context.Background(), "/available_numbers", params)
			if err != nil {
				return fmt.Errorf("could not list available numbers: %w", err)
			}

			for _, item := range items {
				if n, ok := item["number"]; ok {
					fmt.Fprintln(cmd.OutOrStdout(), n)
				}
			}
			return nil
		},
	}
	cmd.Flags().String("country-code", "", "ISO country code (e.g. US)")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	_ = cmd.MarkFlagRequired("country-code")
	return cmd
}
