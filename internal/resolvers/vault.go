package resolvers

import (
	"context"
	"fmt"
	"path"

	"github.com/deref/exo/internal/gensym"
)

type VaultResolver struct {
	Q *RootResolver
	VaultRow
}

type VaultRow struct {
	ID      string `db:"id"`
	StackID string `db:"stack_id"`
	URL     string `db:"url"`
}

type StackVaultRow struct {
	StackID string `db:"stack_id"`
	VaultID string `db:"vault_id"`
}

func (r *QueryResolver) AllVaults(ctx context.Context) ([]*VaultResolver, error) {
	var rows []VaultRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT *
		FROM vault
		ORDER BY url ASC
	`)
	if err != nil {
		return nil, err
	}
	return vaultRowsToResolvers(r, rows), nil
}

func (r *QueryResolver) VaultByID(ctx context.Context, args struct {
	ID string
}) (*VaultResolver, error) {
	return r.vaultByID(ctx, &args.ID)
}

func (r *QueryResolver) vaultByID(ctx context.Context, id *string) (*VaultResolver, error) {
	vault := &VaultResolver{}
	err := r.getRowByKey(ctx, &vault.VaultRow, `
		SELECT *
		FROM vault
		WHERE id = ?
	`, id)
	if vault.ID == "" {
		vault = nil
	}
	return vault, err
}

func (r *QueryResolver) vaultsByStackID(ctx context.Context, stackID string) ([]*VaultResolver, error) {
	var rows []VaultRow
	err := r.db.SelectContext(ctx, &rows, `
		SELECT vault.*
		FROM vault, stack_vault
		WHERE vault.id = stack_vault.vault_id
		AND stack_vault.stack_id = ?
		ORDER BY vault.url ASC
	`, stackID)
	if err != nil {
		return nil, err
	}
	return vaultRowsToResolvers(r, rows), nil
}

func vaultRowsToResolvers(r *RootResolver, rows []VaultRow) []*VaultResolver {
	resolvers := make([]*VaultResolver, len(rows))
	for i, row := range rows {
		resolvers[i] = &VaultResolver{
			Q:        r,
			VaultRow: row,
		}
	}
	return resolvers
}

func (r *MutationResolver) AttachVault(ctx context.Context, args struct {
	StackID string
	URL     string
}) (*VaultResolver, error) {
	row := VaultRow{
		ID:  gensym.RandomBase32(),
		URL: args.URL,
	}
	if err := r.insertRowEx(ctx, "vault", &row, `
		ON CONFLICT ( url )
		DO UPDATE url = url
	`); err != nil {
		return nil, fmt.Errorf("adding vault: %w", err)
	}
	if err := r.insertRow(ctx, "vault_stack", StackVaultRow{}); err != nil {
		return nil, err
	}
	// TODO: Trigger manifest update.
	return &VaultResolver{
		Q:        r,
		VaultRow: row,
	}, nil
}

func (r *MutationResolver) ForgetVault(ctx context.Context, args struct {
	ID string
}) (*VoidResolver, error) {
	_, err := r.db.ExecContext(ctx, `
		BEGIN;

		DELETE FROM vault
		WHERE id = ?;

		DELETE FROM stack_vault
		WHERE vault_id = ?;
		
		COMMIT;
	`, args.ID, args.ID)
	// TODO: Trigger manifest update.
	return nil, err
}

func (r *VaultResolver) Name() string {
	return path.Base(r.URL) // XXX consider scope/uniqueness and editing.
}

func (r *VaultResolver) Error() *string {
	return nil // XXX use status of last cached lookup of secrets, etc.
}

func (r *VaultResolver) Connected() bool {
	return r.Error() == nil // XXX and is not auth error.
}

func (r *VaultResolver) Authenticated() bool {
	return r.Error() == nil
}

func (r *VaultResolver) Secrets(ctx context.Context) ([]*SecretResolver, error) {
	return r.Q.secretsByVaultID(ctx, r.ID)
}
