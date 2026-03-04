package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/auth"
	"github.com/voipbin/vn-cli/internal/output"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

func newTransfersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfers",
		Short: "Manage transfers",
	}
	cmd.AddCommand(
		newTransfersCreateCmd(),
	)
	return cmd
}

var transferDetailColumns = []output.Column{
	{Name: "ID", Field: "id"},
	{Name: "CUSTOMER_ID", Field: "customer_id"},
	{Name: "TYPE", Field: "type"},
	{Name: "TRANSFERER_CALL_ID", Field: "transferer_call_id"},
	{Name: "TRANSFEREE_CALL_ID", Field: "transferee_call_id"},
	{Name: "GROUPCALL_ID", Field: "groupcall_id"},
	{Name: "CREATED", Field: "tm_create"},
}

func newTransfersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := auth.NewClientFromContext(cmd)
			if err != nil {
				return err
			}

			transfererCallID, _ := cmd.Flags().GetString("transferer-call-id")
			transfereeAddr, _ := cmd.Flags().GetString("transferee-address")
			transferType, _ := cmd.Flags().GetString("type")

			body := voipbin_client.PostTransfersJSONRequestBody{
				TransfererCallId:    transfererCallID,
				TransfereeAddresses: []voipbin_client.CommonAddress{{Target: &transfereeAddr}},
				TransferType:        voipbin_client.TransferManagerTransferType(transferType),
			}

			resp, err := client.PostTransfersWithResponse(context.Background(), body)
			if err != nil {
				return fmt.Errorf("could not create transfer: %w", err)
			}
			if resp.StatusCode() != 200 {
				return fmt.Errorf("API error: %s", resp.Status())
			}
			if resp.JSON200 == nil {
				return fmt.Errorf("unexpected empty response")
			}

			return output.PrintItem(cmd, resp.JSON200, transferDetailColumns)
		},
	}
	cmd.Flags().String("transferer-call-id", "", "Call ID of the transferer")
	cmd.Flags().String("transferee-address", "", "Address of the transferee")
	cmd.Flags().String("type", "", "Transfer type")
	_ = cmd.MarkFlagRequired("transferer-call-id")
	_ = cmd.MarkFlagRequired("transferee-address")
	_ = cmd.MarkFlagRequired("type")
	return cmd
}
