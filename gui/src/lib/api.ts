import { isRunning } from './global/server-status';
import type { GetVersionResponse } from './kernel/types';
import type {
  ExportProcfileResponse,
  LogsResponse,
  ReadFileResponse,
} from './logs/types';
import type {
  CreateProcessResponse,
  ProcessDescription,
} from './process/types';
import type { VolumeDescription, NetworkDescription } from './docker/types';
import type { TaskDescription } from './tasks/types';

interface IdleRequest {
  stage: 'idle';
}

interface PendingRequest {
  stage: 'pending';
}

interface ErrorResponse {
  stage: 'error';
  message: string;
}

interface SuccessResponse<T> {
  stage: 'success';
  data: T;
}

interface RefetchingResponse<T> {
  stage: 'refetching';
  data: T;
}

export type RequestLifecycle<T> =
  | IdleRequest
  | PendingRequest
  | ErrorResponse
  | SuccessResponse<T>
  | RefetchingResponse<T>;

export const notRequested = <T>(): RequestLifecycle<T> => ({ stage: 'idle' });
export const pendingRequest = <T>(): RequestLifecycle<T> => ({
  stage: 'pending',
});
export const errorResponse = <T>(message: string): RequestLifecycle<T> => ({
  stage: 'error',
  message,
});
export const successResponse = <T>(data: T): RequestLifecycle<T> => ({
  stage: 'success',
  data,
});
export const refetchingResponse = <T>(prev: T): RequestLifecycle<T> => ({
  stage: 'refetching',
  data: prev,
});

type HasData<T> = SuccessResponse<T> | RefetchingResponse<T>;
// TODO: Should idle be considered unresolved?
type IsUnresolved<T> = IdleRequest | PendingRequest | RefetchingResponse<T>;
type IsResolved<T> = ErrorResponse | SuccessResponse<T>;

export const hasData = <T>(r: RequestLifecycle<T>): r is HasData<T> =>
  r.stage === 'success' || r.stage === 'refetching';

export const isUnresolved = <T>(r: RequestLifecycle<T>): r is IsUnresolved<T> =>
  r.stage === 'idle' || r.stage === 'pending' || r.stage === 'refetching';

export const isResolved = <T>(r: RequestLifecycle<T>): r is IsResolved<T> =>
  r.stage === 'error' || r.stage === 'success';

export interface PaginationParams {
  cursor: string | null;
  prev?: number;
  next?: number;
}

const apiPort = parseInt(import.meta.env.VITE_API_PORT as string);
const baseUrl = `http://localhost:${apiPort}/_exo`;

const apiUrl = (path: string, query: Record<string, string>) => {
  let qs = '';
  let sep = '?';
  for (const [key, value] of Object.entries(query)) {
    qs += sep;
    sep = '&';
    qs += encodeURIComponent(key);
    qs += '=';
    qs += encodeURIComponent(value);
  }
  return baseUrl + path + qs;
};

const isErrorLike = (x: any): x is { message: string } => {
  return (
    x != null &&
    typeof x === 'object' &&
    'message' in x &&
    typeof x.message === 'string'
  );
};

export class APIError extends Error {
  constructor(public readonly httpStatus: number, message: string) {
    super(message);
  }
}

export const isClientError = (err: Error): err is APIError =>
  err instanceof APIError && 400 <= err.httpStatus && err.httpStatus < 500;

const responseToError = async (res: Response): Promise<Error | null> => {
  if (200 <= res.status && res.status < 300) {
    return null;
  }
  const text = await res.text();
  let json: unknown;
  try {
    json = JSON.parse(text);
  } catch (_: unknown) {
    json = text;
  }
  if (!isErrorLike(json)) {
    return new Error(`malformed error from server: ${JSON.stringify(json)}`);
  }
  return new APIError(res.status, json.message);
};

const rpc = async (
  path: string,
  query: Record<string, string>,
  data?: unknown,
): Promise<unknown> => {
  let res: Response;
  try {
    res = await fetch(apiUrl(path, query), {
      method: 'POST',
      headers: {
        accept: 'application/json',
        'content-type': 'application/json',
      },
      ...(data
        ? {
            body: JSON.stringify(data),
          }
        : {}),
    });
  } catch (err: unknown) {
    if (err instanceof TypeError) {
      isRunning.set(false);
      throw new APIError(0, 'Server not available');
    }
    throw err;
  }

  isRunning.update(() => true);

  const err = await responseToError(res);
  if (err !== null) {
    throw err;
  }
  return await res.json();
};

export interface ProcessSpec {
  directory?: string;
  program: string;
  arguments: string[];
  environment?: Record<string, string>;
}

export interface WorkspaceDescription {
  id: string;
  root: string;
}

export interface ComponentDescription {
  id: string;
  name: string;
  type: string;
}

export interface CreateComponentResponse {
  id: string;
}

export interface DescribeTasksInput {
  jobIds?: string[];
}

export interface DirectoryEntry {
  name: string;
  isDirectory: boolean;
}

export interface ReadDirResult {
  directory: DirectoryEntry;
  parent: null | DirectoryEntry;
  entries: DirectoryEntry[];
}

export interface KernelApi {
  getUserHomeDir(): Promise<string>;
  readDir(path: string): Promise<ReadDirResult>;
  describeWorkspaces(): Promise<WorkspaceDescription[]>;
  createWorkspace(root: string): Promise<string>;
  getVersion(): Promise<GetVersionResponse>;
  upgrade(): Promise<void>;
  ping(): Promise<void>;
  describeTasks(input?: DescribeTasksInput): Promise<TaskDescription[]>;
}

export interface WorkspaceApi {
  describeComponents(): Promise<ComponentDescription[]>;

  describeEnvironment(): Promise<Record<string, string>>;

  describeProcesses(): Promise<ProcessDescription[]>;

  describeVolumes(): Promise<VolumeDescription[]>;

  describeNetworks(): Promise<NetworkDescription[]>;

  apply(): Promise<void>;

  createProcess(
    name: string,
    spec: ProcessSpec,
  ): Promise<CreateProcessResponse>;

  createComponent(
    name: string,
    type: string,
    spec: string,
  ): Promise<CreateComponentResponse>;

  startProcess(ref: string): Promise<void>;

  stopProcess(ref: string): Promise<void>;

  deleteComponent(ref: string): Promise<void>;

  refreshAllProcesses(): Promise<void>;

  getEvents(
    logs: string[],
    filterStr: string | null,
    pagination?: PaginationParams,
  ): Promise<LogsResponse>;

  exportProcfile(): Promise<string>;

  readFile(filePath: string): Promise<string | null>;

  writeFile(filePath: string, content: string, mode?: number): Promise<void>;
}

export const api = (() => {
  const kernel: KernelApi = (() => {
    const invoke = (method: string, data?: unknown) =>
      rpc(`/kernel/${method}`, {}, data);
    return {
      async getUserHomeDir(): Promise<string> {
        const { path } = (await invoke('get-user-home-dir', {})) as any;
        return path;
      },
      async readDir(path: string): Promise<ReadDirResult> {
        return (await invoke('read-dir', { path })) as any;
      },
      async describeWorkspaces(): Promise<WorkspaceDescription[]> {
        const { workspaces } = (await invoke('describe-workspaces', {})) as any;
        return workspaces;
      },
      async createWorkspace(root: string): Promise<string> {
        const { id } = (await invoke('create-workspace', { root })) as any;
        return id;
      },

      async getVersion(): Promise<GetVersionResponse> {
        return (await invoke('get-version', {})) as any;
      },

      async upgrade(): Promise<void> {
        return (await invoke('upgrade', {})) as any;
      },

      async ping(): Promise<void> {
        await invoke('ping', {});
      },

      async describeTasks(
        input: DescribeTasksInput = {},
      ): Promise<TaskDescription[]> {
        const { tasks } = (await invoke('describe-tasks', {})) as any;
        return tasks as TaskDescription[];
      },
    };
  })();

  const workspace = (id: string): WorkspaceApi => {
    const invoke = (method: string, data?: unknown) =>
      rpc(`/workspace/${method}`, { id }, data);
    return {
      async describeComponents(): Promise<ComponentDescription[]> {
        const { components } = (await invoke('describe-components')) as any;
        return components;
      },

      async describeEnvironment(): Promise<Record<string, string>> {
        const { variables } = (await invoke('describe-environment')) as any;
        return variables;
      },

      async describeProcesses(): Promise<ProcessDescription[]> {
        const { processes } = (await invoke('describe-processes')) as any;
        return processes;
      },

      async describeVolumes(): Promise<VolumeDescription[]> {
        const { volumes } = (await invoke('describe-volumes')) as any;
        return volumes;
      },

      async describeNetworks(): Promise<NetworkDescription[]> {
        const { networks } = (await invoke('describe-networks')) as any;
        return networks;
      },

      async apply() {
        await invoke('apply', {});
      },

      async createProcess(
        name: string,
        spec: ProcessSpec,
      ): Promise<CreateProcessResponse> {
        return (await invoke('create-component', {
          name,
          type: 'process',
          spec: JSON.stringify(spec),
        })) as CreateProcessResponse;
      },

      async createComponent(
        name: string,
        type: string,
        spec: string,
      ): Promise<CreateComponentResponse> {
        return (await invoke('create-component', {
          name,
          type,
          spec,
        })) as CreateComponentResponse;
      },

      // TODO: Add updateComponent once fully working in backend

      async startProcess(ref: string): Promise<void> {
        await invoke('start-components', { refs: [ref] });
      },

      async stopProcess(ref: string): Promise<void> {
        await invoke('stop-components', { refs: [ref] });
      },

      async deleteComponent(ref: string): Promise<void> {
        await invoke('delete-components', { refs: [ref] });
      },

      async refreshAllProcesses(): Promise<void> {
        await invoke('refresh-components');
      },

      async getEvents(
        logs: string[],
        filterStr: string | null,
        pagination?: PaginationParams,
      ): Promise<LogsResponse> {
        return (await invoke('get-events', {
          logs,
          filterStr,
          ...pagination,
        })) as LogsResponse;
      },

      async exportProcfile(): Promise<string> {
        const res = (await invoke('export-procfile')) as ExportProcfileResponse;
        return res.procfile;
      },

      async readFile(path: string): Promise<string | null> {
        const res = (await invoke('read-file', {
          path,
        })) as ReadFileResponse | null;
        return res?.content ?? null;
      },

      async writeFile(
        path: string,
        content: string,
        mode?: number,
      ): Promise<void> {
        await invoke('write-file', { path, content, mode });
      },
    };
  };

  return {
    kernel,
    workspace,
  };
})();
