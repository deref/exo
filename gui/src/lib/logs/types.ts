
export interface LogEvent {
  id: string;
  timestamp: string;
  log: string; // Process name:(out|err).
  message: string;
}

export interface LogsResponse {
  events: LogEvent[];
  cursor: string | null;
}
