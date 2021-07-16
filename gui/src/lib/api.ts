import type { LogEvent, LogsResponse } from './logs/types';
import type { ProcessDescription } from "./process/types";

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

export type RemoteData<T> = IdleRequest | PendingRequest | ErrorResponse | SuccessResponse<T> | RefetchingResponse<T>;

export const notRequested = <T>(): RemoteData<T> => ({ stage: 'idle' });
export const pendingRequest = <T>(): RemoteData<T> => ({ stage: 'pending' });
export const errorResponse = <T>(message: string): RemoteData<T> => ({ stage: 'error', message });
export const successResponse = <T>(data: T): RemoteData<T> => ({ stage: 'success', data });
export const refetchingResponse = <T>(prev: T): RemoteData<T> => ({ stage: 'refetching', data: prev });

type HasData<T> = SuccessResponse<T> | RefetchingResponse<T>;
// TODO: Should idle be considered unresolved?
type IsUnresolved<T> = IdleRequest | PendingRequest | RefetchingResponse<T>;
type IsResolved<T> = ErrorResponse | SuccessResponse<T>;

export const hasData = <T>(r: RemoteData<T>): r is HasData<T> =>
  r.stage === 'success' || r.stage === 'refetching';

export const IsUnresolved = <T>(r: RemoteData<T>): r is IsUnresolved<T> =>
  r.stage === 'idle' || r.stage === 'pending' || r.stage === 'refetching';

export const IsResolved = <T>(r: RemoteData<T>): r is IsResolved<T> =>
  r.stage === 'error' || r.stage === 'success';

export type PaginationParams = {
  type: 'before-cursor';
  cursor: string | null;
} | {
  type: 'after-cursor';
  cursor: string | null;
};

const baseUrl = 'http://localhost:4000/_exo';
const apiUrl = (path: string) => baseUrl + path;
const rpc = async (path: string, data?: unknown): Promise<unknown> => {
  const res = await fetch(apiUrl(path), {
    method: 'POST',
    headers: {
      accept: 'application/json',
      'content-type': 'application/json',
    },
    ...(data ? {
      body: JSON.stringify(data)
    } : {}),
  });
  return await res.json();
}

export const api = (() => {
  return {
    async describeProcesses(): Promise<ProcessDescription[]> {
      const { processes } = await rpc('/describe-processes') as any;
      return processes;
    },

    async startProcess(ref: string): Promise<void> {
      await rpc('/start', { ref });
    },

    async stopProcess(ref: string): Promise<void> {
      await rpc('/stop', { ref });
    },

    async refreshAllProcesses(): Promise<void> {
      await rpc('/refresh');
    },

    async getEvents(logs: string[], pagination?: PaginationParams): Promise<LogsResponse> {
      if (pagination?.type === 'before-cursor') {
        throw new Error("Before cursor not supported.");
      }
      return await rpc('/get-events', {
        logs,
        ...(pagination?.cursor ? { after: pagination?.cursor } : {}),
      }) as any;
    }
  };
})();
