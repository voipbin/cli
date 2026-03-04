package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			pageToken, _ := cmd.Flags().GetString("page-token")
			pageSize, _ := cmd.Flags().GetInt("page-size")

			params := &voipbin_client.GetAicallsParams{}
			if pageToken != "" {
				params.PageToken = &pageToken
			}
			if pageSize > 0 {
				ps := pageSize
				params.PageSize = &ps
			}

			resp, err := client.GetAicallsWithResponse(context.Background(), params)
			if err != nil {
				return fmt.Errorf("could not list AI calls: %w", err)
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

			return output.PrintList(cmd, *resp.JSON200.Result, aicallListColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.GetAicallsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not get AI call: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aicallDetailColumns)
		},
	}
}

func newAicallsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI call",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			aiID, _ := cmd.Flags().GetString("ai-id")
			referenceID, _ := cmd.Flags().GetString("reference-id")
			referenceType, _ := cmd.Flags().GetString("reference-type")
			language, _ := cmd.Flags().GetString("language")
			gender, _ := cmd.Flags().GetString("gender")

			body := voipbin_client.PostAicallsJSONRequestBody{
				AiId:          aiID,
				ReferenceId:   referenceID,
				ReferenceType: voipbin_client.AIManagerAIcallReferenceType(referenceType),
				Language:      language,
				Gender:        voipbin_client.AIManagerAIcallGender(gender),
			}

			resp, err := client.PostAicallsWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create AI call: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, aicallDetailColumns)
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
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			resp, err := client.DeleteAicallsIdWithResponse(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("could not delete AI call: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}

			fmt.Fprintf(cmd.OutOrStdout(), "AI call %s deleted.\n", args[0])
			return nil
		},
	}
}
