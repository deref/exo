export interface ApiGatewayDescription {
  id: string;
  name: string;
  apiPort: number;
  webPort: number;
}

export interface ProcessDescription {
  id: string;
  name: string;
  running: boolean;
  ports: number[];
  envVars: null | Record<string, string>;
  cpuPercent: null | number;
  createTime: null | number;
  residentMemory: null | number;
  childrenExecutables: null | string[];
}

export interface CreateProcessResponse {
  id: string;
}
