export interface LogEvent {
  id: string;
  timestamp: string;
  log: string; // Process name:(out|err).
  message: string;
}

export interface LogsResponse {
  items: LogEvent[];
  prevCursor: string;
  nextCursor: string;
}

export interface ExportProcfileResponse {
  procfile: string;
}
