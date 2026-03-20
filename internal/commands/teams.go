package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
)

func newTeamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "teams",
		Short: "Manage teams",
	}
	cmd.AddCommand(
		newTeamsListCmd(),
		newTeamsGetCmd(),
		newTeamsCreateCmd(),
		newTeamsUpdateCmd(),
		newTeamsDeleteCmd(),
	)
	return cmd
}

var teamListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "START_MEMBER_ID", Field: "start_member_id"},
	{Name: "CREATED", Field: "tm_create"},
}

var teamDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "NAME", Field: "name"},
	{Name: "DETAIL", Field: "detail"},
	{Name: "START_MEMBER_ID", Field: "start_member_id"},
	{Name: "MEMBERS", Field: "members"},
	{Name: "PARAMETER", Field: "parameter"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newTeamsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List teams",
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

			items, nextToken, err := c.List(context.Background(), "/teams", params)
			if err != nil {
				return fmt.Errorf("could not list teams: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, teamListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newTeamsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a team by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/teams/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get team: %w", err)
			}

			return output.PrintItem(cmd, result, teamDetailColumns)
		},
	}
}

func newTeamsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new team",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			startMemberID, _ := cmd.Flags().GetString("start-member-id")
			membersJSON, _ := cmd.Flags().GetString("members")
			parameterJSON, _ := cmd.Flags().GetString("parameter")

			body := map[string]interface{}{
				"name":   name,
				"detail": detail,
			}
			if startMemberID != "" {
				body["start_member_id"] = startMemberID
			}
			if membersJSON != "" {
				var members interface{}
				if err := json.Unmarshal([]byte(membersJSON), &members); err != nil {
					return fmt.Errorf("invalid members JSON: %w", err)
				}
				body["members"] = members
			}
			if parameterJSON != "" {
				var parameter interface{}
				if err := json.Unmarshal([]byte(parameterJSON), &parameter); err != nil {
					return fmt.Errorf("invalid parameter JSON: %w", err)
				}
				body["parameter"] = parameter
			}

			result, err := c.Post(context.Background(), "/teams", body)
			if err != nil {
				return fmt.Errorf("could not create team: %w", err)
			}

			return output.PrintItem(cmd, result, teamDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Team name")
	cmd.Flags().String("detail", "", "Team detail")
	cmd.Flags().String("start-member-id", "", "Starting member ID")
	cmd.Flags().String("members", "", "Members as JSON string")
	cmd.Flags().String("parameter", "", "Parameter as JSON string")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newTeamsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			name, _ := cmd.Flags().GetString("name")
			detail, _ := cmd.Flags().GetString("detail")
			startMemberID, _ := cmd.Flags().GetString("start-member-id")
			membersJSON, _ := cmd.Flags().GetString("members")
			parameterJSON, _ := cmd.Flags().GetString("parameter")

			body := map[string]interface{}{}
			if name != "" {
				body["name"] = name
			}
			if detail != "" {
				body["detail"] = detail
			}
			if startMemberID != "" {
				body["start_member_id"] = startMemberID
			}
			if membersJSON != "" {
				var members interface{}
				if err := json.Unmarshal([]byte(membersJSON), &members); err != nil {
					return fmt.Errorf("invalid members JSON: %w", err)
				}
				body["members"] = members
			}
			if parameterJSON != "" {
				var parameter interface{}
				if err := json.Unmarshal([]byte(parameterJSON), &parameter); err != nil {
					return fmt.Errorf("invalid parameter JSON: %w", err)
				}
				body["parameter"] = parameter
			}

			result, err := c.Put(context.Background(), "/teams/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update team: %w", err)
			}

			return output.PrintItem(cmd, result, teamDetailColumns)
		},
	}
	cmd.Flags().String("name", "", "Team name")
	cmd.Flags().String("detail", "", "Team detail")
	cmd.Flags().String("start-member-id", "", "Starting member ID")
	cmd.Flags().String("members", "", "Members as JSON string")
	cmd.Flags().String("parameter", "", "Parameter as JSON string")
	return cmd
}

func newTeamsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/teams/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete team: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Team %s deleted.\n", args[0])
			return nil
		},
	}
}
