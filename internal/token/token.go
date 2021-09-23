package token

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

type TokenClient interface {
	GetToken() string
	CheckToken(token string) bool
}

type fileTokenClient struct {
	path   string
	tokens []string
}

var _ TokenClient = &fileTokenClient{}

func NewFileTokenClient(path string) (*fileTokenClient, error) {
	tokenFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reading token file: %w", err)
	}

	scanner := bufio.NewScanner(tokenFile)
	tokens := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tokens = append(tokens, line)
		}
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens in file: %q", path)
	}

	return &fileTokenClient{
		path:   path,
		tokens: tokens,
	}, nil
}

func (c *fileTokenClient) GetToken() string {
	return c.tokens[0]
}

func (c *fileTokenClient) CheckToken(token string) bool {
	for _, t := range c.tokens {
		if t == token {
			return true
		}
	}
	return false
}

var tokenLength = 20

func genToken() (string, error) {
	buff := make([]byte, 20)
	if _, err := rand.Read(buff); err != nil {
		return "", fmt.Errorf("getting randomness: %w", err)
	}
	return hex.EncodeToString(buff), nil
}

func EnsureTokenFile(path string) (*fileTokenClient, error) {
	tokenFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if os.IsExist(err) {
		return NewFileTokenClient(path)
	}
	if err != nil {
		return nil, fmt.Errorf("opening token file: %w", err)
	}

	token, err := genToken()
	if err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}

	_, err = tokenFile.WriteString(token + "\n")
	if err != nil {
		return nil, fmt.Errorf("writing token: %w", err)
	}

	if err := tokenFile.Close(); err != nil {
		return nil, fmt.Errorf("closing token file: %w", err)
	}

	return NewFileTokenClient(path)
}
