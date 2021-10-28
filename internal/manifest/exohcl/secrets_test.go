package exohcl

import (
	"context"
	"testing"
)

// TODO: Test comment/whitespace preservation/insertion, etc.

func TestAppendSecrets(t *testing.T) {
	ctx := context.Background()
	assertRewrite(t, AppendSecrets{
		Context: ctx,
		Source:  "https://example.com/secrets/1",
	}, `
exo = "0.1"
`, `
exo = "0.1"
environment {
  secrets {
    source = "https://example.com/secrets/1"
  }
}
`)
	assertRewrite(t, AppendSecrets{
		Context: ctx,
		Source:  "https://example.com/secrets/2",
	}, `
exo = "0.1"
environment {
  secrets {
    source = "https://example.com/secrets/1"
  }
}
`, `
exo = "0.1"
environment {
  secrets {
    source = "https://example.com/secrets/1"
  }
  secrets {
    source = "https://example.com/secrets/2"
  }
}
`)
}

func TestRemoveSecrets(t *testing.T) {
	ctx := context.Background()
	assertRewrite(t, RemoveSecrets{
		Context: ctx,
		Source:  "https://example.com/secrets/1",
	}, `
exo = "0.1"
`, `
exo = "0.1"
`)
	assertRewrite(t, RemoveSecrets{
		Context: ctx,
		Source:  "https://example.com/secrets/1",
	}, `
exo = "0.1"
environment {
  secrets {
    source = "https://example.com/secrets/1"
  }
}
`, `
exo = "0.1"
`)
	assertRewrite(t, RemoveSecrets{
		Context: ctx,
		Source:  "https://example.com/secrets/1",
	}, `
exo = "0.1"
environment {
  secrets {
    source = "https://example.com/secrets/1"
  }
  secrets {
    source = "https://example.com/secrets/2"
  }
}
`, `
exo = "0.1"
environment {
  secrets {
    source = "https://example.com/secrets/2"
  }
}
`)
}
