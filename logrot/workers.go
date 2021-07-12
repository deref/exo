package logrot

type worker struct {
}

// XXX read all the logs, prefix each with a timestamp, do log rotation.

// Log line format is `<timestamp> <sid> <message>`.
// Timestamp is ISO-8601.
// SID is the incrementing sequence id.
// For example:
//   2021-06-11T12:30.123 12345 Something interesting happened.

func (svc *service) startWorker(logName string) {
}

func (svc *service) stopWorker(logName string) {
}
