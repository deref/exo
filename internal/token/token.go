package token

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deref/exo/internal/gensym"
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

func genToken() string {
	return gensym.RandomBase32()
}

func EnsureTokenFile(path string) (*fileTokenClient, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		token := genToken()
		if err := ioutil.WriteFile(path, []byte(token), 0600); err != nil {
			return nil, fmt.Errorf("writing token file: %w", err)
		}
	}
	return NewFileTokenClient(path)
}
