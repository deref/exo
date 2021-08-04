export interface ProcessDescription {
  id: string;
  name: string;
  running: boolean;
  envVars: null | Record<string, string>;
  cpuPercent: number;
  createTime: number;
  residentMemory: number;
  ports: number[];
  childrenExecutables: string[];
}

export interface CreateProcessResponse {
  id: string;
}
