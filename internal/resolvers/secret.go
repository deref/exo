package resolvers

import "context"

type SecretResolver struct {
	Q *RootResolver
	SecretRecord
}

// NOTE: Not a database row; avoiding secrets on disk.
type SecretRecord struct {
	VaultID string
	Name    string
	Value   string
}

func (r *QueryResolver) secretsByVaultID(ctx context.Context, stackID string) ([]*SecretResolver, error) {
	panic("TODO: secretsByVaultID") // Lookup in a secrets cache.
}

func (r *QueryResolver) secretsByStackID(ctx context.Context, stackID string) ([]*SecretResolver, error) {
	panic("TODO: secretsByStackID") // Lookup in a secrets cache.
}

func (r *SecretResolver) Vault(ctx context.Context) (*VaultResolver, error) {
	return r.Q.vaultByID(ctx, &r.VaultID)
}

func (r *SecretResolver) ValueIf(ctx context.Context, args struct {
	Reveal bool
}) *string {
	if args.Reveal {
		return &r.Value
	} else {
		return nil
	}
}
