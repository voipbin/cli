package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/voipbin/vn-cli/internal/config"
	"github.com/voipbin/voipbin-go/gens/voipbin_client"
)

type accessKeyTransport struct {
	accessKey string
}

func (t *accessKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	newURL := *req.URL
	query := newURL.Query()
	query.Set("accesskey", t.accessKey)
	newURL.RawQuery = query.Encode()

	newReq := req.Clone(req.Context())
	newReq.URL = &newURL
	return http.DefaultClient.Do(newReq)
}

func NewClientFromContext(cmd *cobra.Command) (voipbin_client.ClientWithResponsesInterface, error) {
	accessKey, err := resolveAccessKey(cmd)
	if err != nil {
		return nil, err
	}

	apiURL := resolveAPIURL(cmd)

	return voipbin_client.NewClientWithResponses(apiURL, func(c *voipbin_client.Client) error {
		c.Client = &http.Client{
			Transport: &accessKeyTransport{accessKey: accessKey},
		}
		return nil
	})
}

func resolveAPIURL(cmd *cobra.Command) string {
	if u, _ := cmd.Flags().GetString("api-url"); u != "" {
		return u
	}

	cfg, err := config.Load()
	if err != nil {
		return config.DefaultAPIURL
	}

	profileName, _ := cmd.Flags().GetString("profile")
	if profileName == "" {
		profileName = cfg.CurrentProfile
	}

	if p, ok := cfg.GetProfile(profileName); ok && p.APIURL != "" {
		return p.APIURL
	}

	return config.DefaultAPIURL
}

func resolveAccessKey(cmd *cobra.Command) (string, error) {
	if key, _ := cmd.Flags().GetString("access-key"); key != "" {
		return key, nil
	}

	if key := os.Getenv("VN_ACCESS_KEY"); key != "" {
		return key, nil
	}

	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("no access key found: set --access-key, VN_ACCESS_KEY env, or run 'vn login'")
	}

	profileName, _ := cmd.Flags().GetString("profile")
	if profileName == "" {
		profileName = cfg.CurrentProfile
	}

	p, ok := cfg.GetProfile(profileName)
	if !ok || p.AccessKey == "" {
		return "", fmt.Errorf("no access key found for profile %q: run 'vn login' or set VN_ACCESS_KEY", profileName)
	}

	return p.AccessKey, nil
}
