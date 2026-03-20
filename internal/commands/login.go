package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/client"
	"github.com/voipbin/vn-cli/internal/config"
)

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate and store credentials",
		Long:  "Interactively configure an access key and store it in ~/.vn/config.yaml",
		RunE:  runLogin,
	}
	cmd.Flags().String("access-key", "", "Access key (non-interactive)")
	cmd.Flags().String("profile", "", "Profile name (default: \"default\")")
	cmd.Flags().String("api-url", "", "Custom API URL")
	return cmd
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials for a profile",
		RunE:  runLogout,
	}
}

func runLogin(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	profileName, _ := cmd.Flags().GetString("profile")
	if profileName == "" {
		fmt.Print("Profile name [default]: ")
		input, _ := reader.ReadString('\n')
		profileName = strings.TrimSpace(input)
		if profileName == "" {
			profileName = "default"
		}
	}

	accessKey, _ := cmd.Flags().GetString("access-key")
	if accessKey == "" {
		fmt.Print("Access key: ")
		input, _ := reader.ReadString('\n')
		accessKey = strings.TrimSpace(input)
		if accessKey == "" {
			return fmt.Errorf("access key is required")
		}
	}

	apiURL, _ := cmd.Flags().GetString("api-url")
	if apiURL == "" {
		fmt.Printf("API URL [%s]: ", config.DefaultAPIURL)
		input, _ := reader.ReadString('\n')
		apiURL = strings.TrimSpace(input)
	}
	if apiURL == "" {
		apiURL = config.DefaultAPIURL
	}

	// Validate the access key by calling GetMe
	fmt.Print("Validating access key... ")
	c := client.New(apiURL, accessKey)
	_, err := c.Get(context.Background(), "/me")
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	fmt.Println("OK")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	profile := config.Profile{
		AccessKey: accessKey,
	}
	if apiURL != config.DefaultAPIURL {
		profile.APIURL = apiURL
	}

	cfg.SetProfile(profileName, profile)
	cfg.CurrentProfile = profileName

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}

	fmt.Printf("Credentials saved to profile %q.\n", profileName)
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	profileName, _ := cmd.Flags().GetString("profile")
	if profileName == "" {
		profileName = "default"
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	if _, ok := cfg.GetProfile(profileName); !ok {
		return fmt.Errorf("profile %q not found", profileName)
	}

	cfg.DeleteProfile(profileName)
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}

	fmt.Printf("Profile %q removed.\n", profileName)
	return nil
}
