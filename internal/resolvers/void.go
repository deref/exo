package resolvers

type VoidResolver struct {
	// See note: [EMPTY_GRAPHQL_TYPE].
	Void *VoidResolver
}

var Void *VoidResolver = nil
