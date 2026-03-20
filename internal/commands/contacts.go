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

func newContactsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contacts",
		Short: "Manage contacts",
	}
	cmd.AddCommand(
		newContactsListCmd(),
		newContactsGetCmd(),
		newContactsCreateCmd(),
		newContactsUpdateCmd(),
		newContactsDeleteCmd(),
		newContactsLookupCmd(),
	)
	return cmd
}

var contactListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "DISPLAY_NAME", Field: "display_name"},
	{Name: "COMPANY", Field: "company"},
	{Name: "JOB_TITLE", Field: "job_title"},
	{Name: "SOURCE", Field: "source"},
	{Name: "CREATED", Field: "tm_create"},
}

var contactDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "FIRST_NAME", Field: "first_name"},
	{Name: "LAST_NAME", Field: "last_name"},
	{Name: "DISPLAY_NAME", Field: "display_name"},
	{Name: "COMPANY", Field: "company"},
	{Name: "JOB_TITLE", Field: "job_title"},
	{Name: "SOURCE", Field: "source"},
	{Name: "EXTERNAL_ID", Field: "external_id"},
	{Name: "NOTES", Field: "notes"},
	{Name: "PHONE_NUMBERS", Field: "phone_numbers"},
	{Name: "EMAILS", Field: "emails"},
	{Name: "TAG_IDS", Field: "tag_ids"},
	{Name: "CREATED", Field: "tm_create"},
	{Name: "UPDATED", Field: "tm_update"},
}

func newContactsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List contacts",
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

			items, nextToken, err := c.List(context.Background(), "/contacts", params)
			if err != nil {
				return fmt.Errorf("could not list contacts: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, contactListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newContactsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a contact by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/contacts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get contact: %w", err)
			}

			return output.PrintItem(cmd, result, contactDetailColumns)
		},
	}
}

func newContactsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			firstName, _ := cmd.Flags().GetString("first-name")
			lastName, _ := cmd.Flags().GetString("last-name")
			displayName, _ := cmd.Flags().GetString("display-name")
			company, _ := cmd.Flags().GetString("company")
			jobTitle, _ := cmd.Flags().GetString("job-title")
			source, _ := cmd.Flags().GetString("source")
			externalID, _ := cmd.Flags().GetString("external-id")
			notes, _ := cmd.Flags().GetString("notes")

			body := map[string]interface{}{
				"first_name":   firstName,
				"last_name":    lastName,
				"display_name": displayName,
				"company":      company,
				"job_title":    jobTitle,
				"source":       source,
				"external_id":  externalID,
				"notes":        notes,
			}

			result, err := c.Post(context.Background(), "/contacts", body)
			if err != nil {
				return fmt.Errorf("could not create contact: %w", err)
			}

			return output.PrintItem(cmd, result, contactDetailColumns)
		},
	}
	cmd.Flags().String("first-name", "", "First name")
	cmd.Flags().String("last-name", "", "Last name")
	cmd.Flags().String("display-name", "", "Display name")
	cmd.Flags().String("company", "", "Company")
	cmd.Flags().String("job-title", "", "Job title")
	cmd.Flags().String("source", "", "Source")
	cmd.Flags().String("external-id", "", "External ID")
	cmd.Flags().String("notes", "", "Notes")
	return cmd
}

func newContactsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a contact",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			firstName, _ := cmd.Flags().GetString("first-name")
			lastName, _ := cmd.Flags().GetString("last-name")
			displayName, _ := cmd.Flags().GetString("display-name")
			company, _ := cmd.Flags().GetString("company")
			jobTitle, _ := cmd.Flags().GetString("job-title")
			notes, _ := cmd.Flags().GetString("notes")

			body := map[string]interface{}{}
			if firstName != "" {
				body["first_name"] = firstName
			}
			if lastName != "" {
				body["last_name"] = lastName
			}
			if displayName != "" {
				body["display_name"] = displayName
			}
			if company != "" {
				body["company"] = company
			}
			if jobTitle != "" {
				body["job_title"] = jobTitle
			}
			if notes != "" {
				body["notes"] = notes
			}

			result, err := c.Put(context.Background(), "/contacts/"+args[0], body)
			if err != nil {
				return fmt.Errorf("could not update contact: %w", err)
			}

			return output.PrintItem(cmd, result, contactDetailColumns)
		},
	}
	cmd.Flags().String("first-name", "", "First name")
	cmd.Flags().String("last-name", "", "Last name")
	cmd.Flags().String("display-name", "", "Display name")
	cmd.Flags().String("company", "", "Company")
	cmd.Flags().String("job-title", "", "Job title")
	cmd.Flags().String("notes", "", "Notes")
	return cmd
}

func newContactsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a contact",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/contacts/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete contact: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Contact %s deleted.\n", args[0])
			return nil
		},
	}
}

func newContactsLookupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookup",
		Short: "Lookup contacts by phone number or email",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			phone, _ := cmd.Flags().GetString("phone")
			email, _ := cmd.Flags().GetString("email")

			params := url.Values{}
			if phone != "" {
				params.Set("phone", phone)
			}
			if email != "" {
				params.Set("email", email)
			}

			items, _, err := c.List(context.Background(), "/contacts/lookup", params)
			if err != nil {
				return fmt.Errorf("could not lookup contacts: %w", err)
			}

			return output.PrintList(cmd, items, contactListColumns)
		},
	}
	cmd.Flags().String("phone", "", "Phone number to search")
	cmd.Flags().String("email", "", "Email to search")
	return cmd
}
