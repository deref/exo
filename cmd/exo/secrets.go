package main

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

func requestDeviceCode() (deviceCodeResponse, error) {
	url := "https://deref.us.auth0.com/oauth/device/code"

	payloadString := fmt.Sprintf("client_id=%s&scope=profile&audience=cli-client", clientId)
	payload := strings.NewReader(payloadString)
	req, _ := http.NewRequest("POST", url, payload)

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
	uri := "https://deref.us.auth0.com/oauth/token"
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

var secretsCmd = &cobra.Command{
	Use:  "secrets",
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		accessToken, err := ioutil.ReadFile("exo-secrets.json")
		if err != nil {
			if os.IsNotExist(err) {
				codeResponse, err := requestDeviceCode()
				if err != nil {
					return fmt.Errorf("getting device code: %w", err)
				}

				fmt.Println("Open the following URL to authenticate:")
				fmt.Println(codeResponse.VerificationURIComplete)

				tokens, err := requestTokens(codeResponse.DeviceCode, codeResponse.Interval)
				if err != nil {
					return fmt.Errorf("getting tokens: %w", err)
				}

				err = ioutil.WriteFile(cfg.EsvTokenFile, []byte(tokens.AccessToken), 0600)
				if err != nil {
					return fmt.Errorf("writing secrets: %w", err)
				}

				accessToken = []byte(tokens.AccessToken)
			}
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
