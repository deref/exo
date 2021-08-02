package storage

type OID int32

const (
	// Negative numbers reserved for system objects.
	oidBootstrap OID = -1
	oidTable     OID = -2
	oidSchema    OID = -3
	oidIndex     OID = -4
)
