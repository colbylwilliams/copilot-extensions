package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/colbylwilliams/copilot-extensions/pkg/config"
	gogh "github.com/google/go-github/v64/github"
	"github.com/jferrl/go-githubauth"
	"golang.org/x/oauth2"
	// octokit "github.com/octokit/go-sdk/pkg"
)

const (
	publicKeyTypeCopilotAPI string = "copilot_api"
	publicKeyBaseURL        string = "https://api.github.com/meta/public_keys/"
)

func GetInstallationClient(ctx context.Context, cfg *config.Config, installationID int64) (*gogh.Client, error) {
	ts, err := githubauth.NewApplicationTokenSource(cfg.GitHubAppID, cfg.GitHubAppPrivateKey,
		githubauth.WithApplicationTokenExpiration(5*time.Minute),
	)
	// appTokenSource, err := githubauth.NewApplicationTokenSource(appID, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create app token source: %w", err)
	}
	// oauth2.NewClient create a new http.Client that adds an Authorization header with the token.
	// Transport src use oauth2.ReuseTokenSource to reuse the token.
	// The token will be reused until it expires.
	// The token will be refreshed if it's expired.
	httpClient := oauth2.NewClient(ctx, ts)

	client := gogh.NewClient(httpClient)
	if envURL := os.Getenv("GITHUB_API_URL"); envURL != "" {
		client.BaseURL, _ = url.Parse(envURL + "/")
	}
	return client, nil
}

func GetUserClient(cfg *config.Config, apiToken string) *gogh.Client {
	client := gogh.NewClient(nil)
	if envURL := os.Getenv("GITHUB_API_URL"); envURL != "" {
		client.BaseURL, _ = url.Parse(envURL + "/")
	}
	return client.WithAuthToken(apiToken)
}

func FetchPublicKey() (string, error) {

	// only one key type for now
	keyType := publicKeyTypeCopilotAPI

	resp, err := http.Get(publicKeyBaseURL + keyType)
	if err != nil {
		return "", fmt.Errorf("failed to fetch public key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch public key: %s", resp.Status)
	}

	var respBody struct {
		PublicKeys []struct {
			Identifier string `json:"key_identifier"`
			Key        string `json:"key"`
			IsCurrent  bool   `json:"is_current"`
		} `json:"public_keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}

	for _, pk := range respBody.PublicKeys {
		if pk.IsCurrent {
			return pk.Key, nil
		}
	}

	return "", errors.New("could not find current public key")
}
