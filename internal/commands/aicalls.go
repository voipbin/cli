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

func newAicallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aicalls",
		Short: "Manage AI calls",
	}
	cmd.AddCommand(
		newAicallsListCmd(),
		newAicallsGetCmd(),
		newAicallsCreateCmd(),
		newAicallsDeleteCmd(),
	)
	return cmd
}

var aicallListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "AI_ID", Field: "ai_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var aicallDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "AI_ID", Field: "ai_id"},
	{Name: "REFERENCE_ID", Field: "reference_id"},
	{Name: "REFERENCE_TYPE", Field: "reference_type"},
	{Name: "STATUS", Field: "status"},
	{Name: "LANGUAGE", Field: "language"},
	{Name: "GENDER", Field: "gender"},
	{Name: "ENGINE_TYPE", Field: "engine_type"},
	{Name: "ENGINE_MODEL", Field: "engine_model"},
	{Name: "TRANSCRIBE_ID", Field: "transcribe_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newAicallsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List AI calls",
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

			items, nextToken, err := c.List(context.Background(), "/aicalls", params)
			if err != nil {
				return fmt.Errorf("could not list AI calls: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, aicallListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newAicallsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an AI call by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/aicalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get AI call: %w", err)
			}

			return output.PrintItem(cmd, result, aicallDetailColumns)
		},
	}
}

func newAicallsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI call",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			aiID, _ := cmd.Flags().GetString("ai-id")
			referenceID, _ := cmd.Flags().GetString("reference-id")
			referenceType, _ := cmd.Flags().GetString("reference-type")
			language, _ := cmd.Flags().GetString("language")
			gender, _ := cmd.Flags().GetString("gender")

			body := map[string]interface{}{
				"ai_id":          aiID,
				"reference_id":   referenceID,
				"reference_type": referenceType,
				"language":       language,
				"gender":         gender,
			}

			result, err := c.Post(context.Background(), "/aicalls", body)
			if err != nil {
				return fmt.Errorf("could not create AI call: %w", err)
			}

			return output.PrintItem(cmd, result, aicallDetailColumns)
		},
	}
	cmd.Flags().String("ai-id", "", "AI ID")
	cmd.Flags().String("reference-id", "", "Reference ID")
	cmd.Flags().String("reference-type", "", "Reference type")
	cmd.Flags().String("language", "", "Language")
	cmd.Flags().String("gender", "", "Gender")
	_ = cmd.MarkFlagRequired("ai-id")
	_ = cmd.MarkFlagRequired("reference-id")
	_ = cmd.MarkFlagRequired("reference-type")
	return cmd
}

func newAicallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an AI call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/aicalls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI call %s deleted.\n", args[0])
			return nil
		},
	}
}
