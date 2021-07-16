package api

type Provider interface {
	Lifecycle
	Process // TODO: Do not require all providers to implement.
}
