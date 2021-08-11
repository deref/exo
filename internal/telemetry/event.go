package telemetry

import "github.com/deref/exo"

type event interface {
	ID() string
	Payload() map[string]interface{}
}

type eventRequest struct {
	ID      string
	Payload string
}

type SystemInfoIdentified struct {
}

func (e *SystemInfoIdentified) ID() string {
	return "TODO"
}

func (e *SystemInfoIdentified) Payload() map[string]interface{} {
	return map[string]interface{}{
		"some": "data",
		"goes": "here",
	}
}

type OperationsPerformed struct {
	Operation       string
	Success         bool
	DurationSummary SummaryStatistics
}

func (e *OperationsPerformed) ID() string {
	return "a20fc050-ba26-4f05-832a-7d42283d1888"
}

func (e *OperationsPerformed) Payload() map[string]interface{} {
	return map[string]interface{}{
		"exoVersion": exo.Version,
		"operation":  e.Operation,
		"success":    e.Success,
		"summary":    e.DurationSummary,
	}
}
