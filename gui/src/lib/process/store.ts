import { writable } from 'svelte/store';
import { api, errorResponse, notRequested, pendingRequest, refetchingResponse, successResponse } from '../api';
import type { ProcessDescription } from './types';

export const processes = writable(notRequested<ProcessDescription[]>());

const refetchProcesses = () =>
  api.describeProcesses()
    .then((data) => {
      processes.set(successResponse(data));
    })
    .catch((err: Error) => {
      processes.set(errorResponse(err.message));
    });

export const fetchProcesses = () => {
  processes.update((req) => {
    switch (req.stage) {
      case 'idle':
        refetchProcesses();
        return pendingRequest();
      case 'pending':
        // Do not refetch.
        return req;
      case 'error':
        refetchProcesses();
        return pendingRequest();
      case 'success':
        refetchProcesses();
        return refetchingResponse(req.data);
    }
  });
}

export const startProcess = async (id: string) =>
  api.startProcess(id)
    .then(() => fetchProcesses());

export const stopProcess = async (id: string) =>
  api.stopProcess(id)
    .then(() => fetchProcesses());

export const refreshAllProcesses = async () =>
  api.refreshAllProcesses()
  .then(() => fetchProcesses());
