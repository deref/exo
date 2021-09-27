package token

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/deref/exo/internal/gensym"
)

type TokenClient interface {
	GetToken() (string, error)
	CheckToken(token string) (bool, error)
}

func genToken() string {
	return gensym.RandomBase32()
}

type FileTokenClient struct {
	Path string
}

var _ TokenClient = &FileTokenClient{}

func readTokenFile(path string) ([]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading token file: %w", err)
	}
	return strings.Fields(string(data)), nil
}

func (c *FileTokenClient) GetToken() (string, error) {
	tokens, err := readTokenFile(c.Path)
	if err != nil {
		return "", fmt.Errorf("reading token file: %w", err)
	}

	if len(tokens) == 0 {
		return "", errors.New("no token in tokens file")
	}
	return tokens[0], nil
}

func (c *FileTokenClient) CheckToken(tokenToCheck string) (bool, error) {
	if tokenToCheck == "" {
		return false, nil
	}

	tokens, err := readTokenFile(c.Path)
	if err != nil {
		return false, fmt.Errorf("reading token file: %w", err)
	}

	for _, token := range tokens {
		if tokenToCheck == token {
			return true, nil
		}
	}
	return false, nil
}

func EnsureTokenFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		token := genToken()
		if err := ioutil.WriteFile(path, []byte(token), 0600); err != nil {
			return fmt.Errorf("writing token file: %w", err)
		}
	}
	return nil
}
