package commands

import (
	"github.com/spf13/cobra"
)

var version = "dev"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vn",
		Short: "VoIPBIN CLI - manage your VoIPBIN resources",
		Long:  "vn is a command-line interface for the VoIPBIN API.\nManage calls, messages, agents, campaigns, and more.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global persistent flags
	cmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")
	cmd.PersistentFlags().String("profile", "", "Configuration profile to use")
	cmd.PersistentFlags().String("access-key", "", "API access key (overrides config and env)")
	cmd.PersistentFlags().String("api-url", "", "API base URL (overrides config)")

	// Identity/auth commands
	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newLogoutCmd())
	cmd.AddCommand(newMeCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newCompletionCmd())

	// Resource commands — registered here so they appear in help
	// Communication
	cmd.AddCommand(newCallsCmd())
	cmd.AddCommand(newMessagesCmd())
	cmd.AddCommand(newEmailsCmd())
	cmd.AddCommand(newConferencesCmd())
	cmd.AddCommand(newConferencecallsCmd())
	cmd.AddCommand(newGroupcallsCmd())
	cmd.AddCommand(newTransfersCmd())

	// AI
	cmd.AddCommand(newAisCmd())
	cmd.AddCommand(newAicallsCmd())
	cmd.AddCommand(newAimessagesCmd())
	cmd.AddCommand(newAisummariesCmd())

	// Campaign/Outbound
	cmd.AddCommand(newCampaignsCmd())
	cmd.AddCommand(newCampaigncallsCmd())
	cmd.AddCommand(newOutdialsCmd())
	cmd.AddCommand(newOutplansCmd())

	// Routing/Flow
	cmd.AddCommand(newFlowsCmd())
	cmd.AddCommand(newActiveflowsCmd())
	cmd.AddCommand(newRoutesCmd())
	cmd.AddCommand(newQueuesCmd())
	cmd.AddCommand(newQueuecallsCmd())
	cmd.AddCommand(newExtensionsCmd())

	// Chat/Conversation
	cmd.AddCommand(newChatsCmd())
	cmd.AddCommand(newChatroomsCmd())
	cmd.AddCommand(newChatmessagesCmd())
	cmd.AddCommand(newChatroommessagesCmd())
	cmd.AddCommand(newConversationsCmd())
	cmd.AddCommand(newConversationAccountsCmd())

	// Account/Management
	cmd.AddCommand(newAgentsCmd())
	cmd.AddCommand(newCustomersCmd())
	cmd.AddCommand(newAccesskeysCmd())
	cmd.AddCommand(newTagsCmd())

	// Telecom
	cmd.AddCommand(newNumbersCmd())
	cmd.AddCommand(newAvailableNumbersCmd())
	cmd.AddCommand(newProvidersCmd())
	cmd.AddCommand(newTrunksCmd())

	// Media/Storage
	cmd.AddCommand(newRecordingsCmd())
	cmd.AddCommand(newTranscribesCmd())
	cmd.AddCommand(newFilesCmd())
	cmd.AddCommand(newStorageAccountsCmd())
	cmd.AddCommand(newStorageFilesCmd())

	// Billing
	cmd.AddCommand(newBillingAccountsCmd())
	cmd.AddCommand(newBillingsCmd())

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("vn version", version)
		},
	}
}

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			default:
				return cmd.Help()
			}
		},
	}
	return cmd
}
