package resolvers

import "context"

type ResourceResolver struct {
	Q *QueryResolver
	ResourceRow
}

type ResourceRow struct {
	IRI         string `db:"iri"`
	ComponentID string `db:"componentID"`
}

func (r *QueryResolver) ResourceByIRI(ctx context.Context, args struct {
	IRI string
}) (*ResourceResolver, error) {
	return r.stackByIRI(ctx, &args.IRI)
}

func (r *QueryResolver) stackByIRI(ctx context.Context, iri *string) (*ResourceResolver, error) {
	s := &ResourceResolver{}
	err := r.getRowByID(ctx, &s.ResourceRow, `
		SELECT iri
		FROM resource
		WHERE iri = ?
	`, iri)
	if s.IRI == "" {
		s = nil
	}
	return s, err
}
