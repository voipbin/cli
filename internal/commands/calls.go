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

func newCallsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calls",
		Short: "Manage calls",
	}
	cmd.AddCommand(
		newCallsListCmd(),
		newCallsGetCmd(),
		newCallsCreateCmd(),
		newCallsDeleteCmd(),
		newCallsHangupCmd(),
		newCallsHoldCmd(),
		newCallsUnholdCmd(),
		newCallsMuteCmd(),
		newCallsUnmuteCmd(),
		newCallsSilenceCmd(),
		newCallsUnsilenceCmd(),
		newCallsMohCmd(),
		newCallsUnmohCmd(),
		newCallsRecordingStartCmd(),
		newCallsRecordingStopCmd(),
		newCallsTalkCmd(),
	)
	return cmd
}

var callListColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "SOURCE", Field: "source"},
	{Name: "DESTINATION", Field: "destination"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "STATUS", Field: "status"},
	{Name: "CREATED", Field: "tm_create"},
}

var callDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "OWNER_ID", Field: "owner_id"},
	{Name: "SOURCE", Field: "source"},
	{Name: "DESTINATION", Field: "destination"},
	{Name: "DIRECTION", Field: "direction"},
	{Name: "STATUS", Field: "status"},
	{Name: "FLOW_ID", Field: "flow_id"},
	{Name: "ACTIVEFLOW_ID", Field: "activeflow_id"},
	{Name: "RECORDING_ID", Field: "recording_id"},
	{Name: "GROUPCALL_ID", Field: "groupcall_id"},
	{Name: "HANGUP_BY", Field: "hangup_by"},
	{Name: "HANGUP_REASON", Field: "hangup_reason"},
	{Name: "CREATED", Field: "tm_create"},
}

func newCallsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List calls",
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

			items, nextToken, err := c.List(context.Background(), "/calls", params)
			if err != nil {
				return fmt.Errorf("could not list calls: %w", err)
			}

			if nextToken != "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "Next page token: %s\n", nextToken)
			}

			return output.PrintList(cmd, items, callListColumns)
		},
	}
	cmd.Flags().String("page-token", "", "Pagination token")
	cmd.Flags().Int("page-size", 0, "Number of results per page")
	return cmd
}

func newCallsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a call by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			result, err := c.Get(context.Background(), "/calls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not get call: %w", err)
			}

			return output.PrintItem(cmd, result, callDetailColumns)
		},
	}
}

func newCallsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new call",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			source, _ := cmd.Flags().GetString("source")
			destination, _ := cmd.Flags().GetString("destination")
			flowID, _ := cmd.Flags().GetString("flow-id")

			body := map[string]interface{}{
				"source":       map[string]interface{}{"target": source},
				"destinations": []map[string]interface{}{{"target": destination}},
			}
			if flowID != "" {
				body["flow_id"] = flowID
			}

			result, err := c.Post(context.Background(), "/calls", body)
			if err != nil {
				return fmt.Errorf("could not create call: %w", err)
			}

			return output.PrintItem(cmd, result, callDetailColumns)
		},
	}
	cmd.Flags().String("source", "", "Source address")
	cmd.Flags().String("destination", "", "Destination address")
	cmd.Flags().String("flow-id", "", "Flow ID to execute")
	_ = cmd.MarkFlagRequired("source")
	_ = cmd.MarkFlagRequired("destination")
	return cmd
}

func newCallsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/calls/"+args[0])
			if err != nil {
				return fmt.Errorf("could not delete call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s deleted.\n", args[0])
			return nil
		},
	}
}

func newCallsHangupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hangup <id>",
		Short: "Hang up a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/hangup", nil)
			if err != nil {
				return fmt.Errorf("could not hangup call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s hung up.\n", args[0])
			return nil
		},
	}
}

func newCallsHoldCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hold <id>",
		Short: "Put a call on hold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/hold", nil)
			if err != nil {
				return fmt.Errorf("could not hold call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s on hold.\n", args[0])
			return nil
		},
	}
}

func newCallsUnholdCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unhold <id>",
		Short: "Resume a held call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/calls/"+args[0]+"/hold")
			if err != nil {
				return fmt.Errorf("could not unhold call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s resumed.\n", args[0])
			return nil
		},
	}
}

func newCallsMuteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mute <id>",
		Short: "Mute a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/mute", map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("could not mute call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s muted.\n", args[0])
			return nil
		},
	}
}

func newCallsUnmuteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unmute <id>",
		Short: "Unmute a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/calls/"+args[0]+"/mute")
			if err != nil {
				return fmt.Errorf("could not unmute call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s unmuted.\n", args[0])
			return nil
		},
	}
}

func newCallsSilenceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "silence <id>",
		Short: "Silence a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/silence", nil)
			if err != nil {
				return fmt.Errorf("could not silence call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s silenced.\n", args[0])
			return nil
		},
	}
}

func newCallsUnsilenceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unsilence <id>",
		Short: "Unsilence a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/calls/"+args[0]+"/silence")
			if err != nil {
				return fmt.Errorf("could not unsilence call: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s unsilenced.\n", args[0])
			return nil
		},
	}
}

func newCallsMohCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "moh <id>",
		Short: "Start music on hold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/moh", nil)
			if err != nil {
				return fmt.Errorf("could not start MOH: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s MOH started.\n", args[0])
			return nil
		},
	}
}

func newCallsUnmohCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unmoh <id>",
		Short: "Stop music on hold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Delete(context.Background(), "/calls/"+args[0]+"/moh")
			if err != nil {
				return fmt.Errorf("could not stop MOH: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s MOH stopped.\n", args[0])
			return nil
		},
	}
}

func newCallsRecordingStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recording-start <id>",
		Short: "Start recording a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/recording_start", map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("could not start recording: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s recording started.\n", args[0])
			return nil
		},
	}
}

func newCallsRecordingStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recording-stop <id>",
		Short: "Stop recording a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/recording_stop", nil)
			if err != nil {
				return fmt.Errorf("could not stop recording: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s recording stopped.\n", args[0])
			return nil
		},
	}
}

func newCallsTalkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "talk <id>",
		Short: "Send talk command to a call",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			_, err = c.Post(context.Background(), "/calls/"+args[0]+"/talk", map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("could not send talk: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Call %s talk sent.\n", args[0])
			return nil
		},
	}
}
