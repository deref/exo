export interface GetVersionResponse {
  installed: string;
  latest?: string;
  current: string;
  managed: boolean;
}
