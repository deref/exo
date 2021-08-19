export interface LogEvent {
  id: string;
  style: string;
  time: {
    short: string;
    full: string;
  };
  name: string;
  message: string;
}
