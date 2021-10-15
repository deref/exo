package esv

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/deref/exo/internal/util/logging"
)

var AuthError = errors.New("auth error")

type EsvClient struct {
	TokenPath    string
	refreshToken string
}

type AuthResponse struct {
	UserCode string
	AuthURL  string
}

func (c *EsvClient) StartAuthFlow(ctx context.Context) (AuthResponse, error) {
	codeResponse, err := requestDeviceCode()
	if err != nil {
		return AuthResponse{}, fmt.Errorf("requesting device code: %w", err)
	}

	go func() {
		logger := logging.CurrentLogger(ctx)

		tokens, err := requestTokens(codeResponse.DeviceCode, codeResponse.Interval)
		if err != nil {
			logger.Infof("got error requesting tokens: %s", err)
			return
		}

		c.refreshToken = tokens.RefreshToken

		err = ioutil.WriteFile(c.TokenPath, []byte(tokens.RefreshToken), 0600)
		if err != nil {
			logger.Infof("writing esv secret: %s", err)
			return
		}
	}()

	return AuthResponse{
		AuthURL:  codeResponse.VerificationURIComplete,
		UserCode: codeResponse.UserCode,
	}, nil
}

func (c *EsvClient) getAccessToken() (string, error) {
	if c.TokenPath == "" {
		return "", fmt.Errorf("token file not set")
	}
	if c.refreshToken == "" {
		tokenBytes, err := ioutil.ReadFile(c.TokenPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("%w: token file does not exist", AuthError)
			}
			return "", fmt.Errorf("reading token file: %w", err)
		}
		c.refreshToken = strings.TrimSpace(string(tokenBytes))
	}

	// FIXME: don't refresh access token on every request.
	accessToken, err := getNewAccessToken(c.refreshToken)
	if err != nil {
		return "", fmt.Errorf("getting access token: %w", err)
	}
	return accessToken, nil
}

func (c *EsvClient) runCommand(output interface{}, host, commandName string, body interface{}) error {
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling command body: %w", err)
	}

	accessToken, err := c.getAccessToken()
	if err != nil {
		return fmt.Errorf("getting access token: %w", err)
	}

	req, _ := http.NewRequest("POST", host+"/api/_exo/"+commandName, bytes.NewBuffer(marshalledBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("making token request: %w", err)
	}

	if resp.StatusCode == 401 {
		return fmt.Errorf("running command %q: %w", commandName, AuthError)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading command result: %w", err)
	}

	if err := json.Unmarshal(result, output); err != nil {
		return fmt.Errorf("unmarshalling command result: %w", err)
	}

	return nil

}

func (c *EsvClient) GetWorkspaceSecrets(vaultURL string) (map[string]string, error) {
	type describeVaultResp struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		Secrets     []struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			Value       string `json:"value"`
		} `json:"secrets"`
	}

	organizationID, vaultID, err := getIdsFromUrl(vaultURL)
	if err != nil {
		return nil, fmt.Errorf("could not find IDs: %w", err)
	}

	uri, err := url.Parse(vaultURL)
	if err != nil {
		return nil, fmt.Errorf("parsing secrets URL: %w", err)
	}
	host := url.URL{Scheme: uri.Scheme, Host: uri.Host}

	resp := describeVaultResp{}
	err = c.runCommand(&resp, host.String(), "describe-project", map[string]string{
		"organizationId": organizationID,
		"vaultId":        vaultID,
	})
	if err != nil {
		return nil, fmt.Errorf("running describe-project command: %w", err)
	}
	secrets := map[string]string{}
	for _, secret := range resp.Secrets {
		secrets[secret.DisplayName] = secret.Value
	}
	return secrets, nil
}

func getIdsFromUrl(vaultURL string) (organizationID, vaultID string, err error) {
	parsedUrl, err := url.Parse(vaultURL)
	if err != nil {
		return "", "", fmt.Errorf("parsing vault URL: %w", err)
	}

	parts := strings.Split(parsedUrl.Path, "/")
	for i, part := range parts {
		if part == "organizations" {
			if i+1 < len(parts) {
				organizationID = parts[i+1]
			}
		}
		if part == "vaults" {
			if i+1 < len(parts) {
				vaultID = parts[i+1]
			}
		}
		if organizationID != "" && vaultID != "" {
			return
		}
	}
	err = fmt.Errorf("could not find IDs in URL: %q", vaultURL)
	return
}