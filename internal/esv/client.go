package esv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type EsvClient struct {
	AccessKey string
}

func (c *EsvClient) runCommand(output interface{}, host, commandName string, body interface{}) error {
	marshalledBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshalling command body: %w", err)
	}

	req, _ := http.NewRequest("POST", host+"/api/_exo/"+commandName, bytes.NewBuffer(marshalledBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("making token request: %w", err)
	}

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading command result: %w", err)
	}

	if err := json.Unmarshal([]byte(result), output); err != nil {
		return fmt.Errorf("unmarshalling command result: %w", err)
	}

	return nil

}

func (c *EsvClient) GetWorkspaceSecrets(projectURL string) (map[string]string, error) {
	type describeProjectResp struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		Secrets     []struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"secrets"`
	}

	organizationID, projectID, err := getIdsFromUrl(projectURL)
	if err != nil {
		return nil, fmt.Errorf("could not find IDs: %w", err)
	}

	uri, err := url.Parse(projectURL)
	if err != nil {
		return nil, fmt.Errorf("parsing secrets URL: %w", err)
	}
	host := url.URL{Scheme: uri.Scheme, Host: uri.Host}

	resp := describeProjectResp{}
	err = c.runCommand(&resp, host.String(), "describe-project", map[string]string{
		"organizationId": organizationID,
		"projectId":      projectID,
	})
	if err != nil {
		return nil, fmt.Errorf("running describe-project command: %w", err)
	}
	secrets := map[string]string{}
	for _, secret := range resp.Secrets {
		secrets[secret.DisplayName] = secret.ID
	}
	return secrets, nil
}

func getIdsFromUrl(projectURL string) (organizationID, projectID string, err error) {
	parsedUrl, err := url.Parse(projectURL)
	if err != nil {
		return "", "", fmt.Errorf("parsing project URL: %w", err)
	}

	parts := strings.Split(parsedUrl.Path, "/")
	for i, part := range parts {
		if part == "organizations" {
			if i+1 < len(parts) {
				organizationID = parts[i+1]
			}
		}
		if part == "projects" {
			if i+1 < len(parts) {
				projectID = parts[i+1]
			}
		}
		if organizationID != "" && projectID != "" {
			return
		}
	}
	err = fmt.Errorf("could not find IDs in URL: %q", projectURL)
	return
}
