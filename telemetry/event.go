package telemetry

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
