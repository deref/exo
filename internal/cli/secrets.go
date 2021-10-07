package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(secretsCmd)
}

type deviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

var clientId = "LNPi71pWh6trIbZOGxxGi5eilI5DakWE"
var derefAuth0Domain = "https://deref.us.auth0.com"

func requestDeviceCode() (deviceCodeResponse, error) {
	uri := derefAuth0Domain + "/oauth/device/code"

	values := url.Values{}
	values.Set("audience", "cli-client")
	values.Set("scope", "profile offline_access")
	values.Set("client_id", clientId)
	payload := strings.NewReader(values.Encode())
	req, _ := http.NewRequest("POST", uri, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return deviceCodeResponse{}, fmt.Errorf("getting auth URL: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return deviceCodeResponse{}, fmt.Errorf("unexpected %d status code when getting auth URL", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return deviceCodeResponse{}, fmt.Errorf("reading body: %w", err)
	}

	codeResult := deviceCodeResponse{}
	if err := json.Unmarshal(body, &codeResult); err != nil {
		return deviceCodeResponse{}, fmt.Errorf("unmarshalling auth URL response: %w", err)
	}

	return codeResult, nil
}

func requestTokens(deviceCode string, interval int) (tokenResponse, error) {
	uri := derefAuth0Domain + "/oauth/token"
	values := url.Values{}
	values.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	values.Set("device_code", deviceCode)
	values.Set("client_id", clientId)
	payloadString := values.Encode()

	for {
		// A new request object must be made for each request made.
		req, err := http.NewRequest("POST", uri, strings.NewReader(payloadString))
		if err != nil {
			return tokenResponse{}, fmt.Errorf("making request: %w", err)
		}

		req.Header.Add("content-type", "application/x-www-form-urlencoded")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return tokenResponse{}, fmt.Errorf("making token request: %w", err)
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return tokenResponse{}, fmt.Errorf("reading body: %w", err)
		}

		if res.StatusCode == 200 {
			tokens := tokenResponse{}
			if err := json.Unmarshal(body, &tokens); err != nil {
				return tokenResponse{}, fmt.Errorf("unmarshalling auth token response: %w", err)
			}
			return tokens, nil
		}

		type errResponseStruct struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}

		errResp := errResponseStruct{}
		if err := json.Unmarshal(body, &errResp); err != nil {
			return tokenResponse{}, fmt.Errorf("unmarshalling auth token error response: %w", err)
		}

		if errResp.Error != "authorization_pending" {
			return tokenResponse{}, fmt.Errorf("unexpected error %s: %s", errResp.Error, errResp.ErrorDescription)
		}

		time.Sleep(time.Second * time.Duration(interval))
	}
}

func getNewAccessToken(refreshToken string) (string, error) {
	uri := derefAuth0Domain + "/oauth/token"

	values := url.Values{}
	values.Set("client_id", clientId)
	values.Set("refresh_token", refreshToken)
	values.Set("grant_type", "refresh_token")
	payloadString := values.Encode()

	req, err := http.NewRequest("POST", uri, strings.NewReader(payloadString))
	if err != nil {
		return "", fmt.Errorf("building refresh token request: %w", err)
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("performing refresh token request: %w", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode == 200 {
		tokens := tokenResponse{}
		if err := json.Unmarshal(body, &tokens); err != nil {
			return "", fmt.Errorf("unmarshalling auth token response: %w", err)
		}
		return tokens.AccessToken, nil
	}

	return "", fmt.Errorf("unexpected status: %q", res.Status)
}

var secretsCmd = &cobra.Command{
	Use:  "secrets",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		refreshToken, err := ioutil.ReadFile("exo-secrets.json")
		if err != nil {
			if os.IsNotExist(err) {
				codeResponse, err := requestDeviceCode()
				if err != nil {
					return fmt.Errorf("getting device code: %w", err)
				}

				// This link is opened automatically because it's single use only. That
				// means if we open it in the wrong browser it becomes worthless.
				fmt.Println("Open the following URL to authenticate:")
				fmt.Println(codeResponse.VerificationURIComplete)

				tokens, err := requestTokens(codeResponse.DeviceCode, codeResponse.Interval)
				if err != nil {
					return fmt.Errorf("getting tokens: %w", err)
				}

				refreshToken = []byte(tokens.RefreshToken)
				err = ioutil.WriteFile(cfg.EsvTokenFile, []byte(tokens.RefreshToken), 0600)
				if err != nil {
					return fmt.Errorf("writing secrets: %w", err)
				}
			}
		}

		accessToken, err := getNewAccessToken(string(refreshToken))
		if err != nil {
			return fmt.Errorf("getting refresh token: %w", err)
		}

		req, _ := http.NewRequest("POST", "http://localhost:5000/api/_exo/describe-me", bytes.NewBufferString(""))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+string(accessToken))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("making token request: %w", err)

		}

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("body: %+v\n", string(body))

		return nil
	},
}
