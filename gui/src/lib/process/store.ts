import { writable } from 'svelte/store';
import { api, errorResponse, notRequested, pendingRequest, refetchingResponse, successResponse } from '../api';
import type { ProcessDescription } from './types';

export const processes = writable(notRequested<ProcessDescription[]>());

const refetchProcesses = (workspace) =>
  workspace.describeProcesses()
    .then((data) => {
      processes.set(successResponse(data));
    })
    .catch((err: Error) => {
      processes.set(errorResponse(err.message));
    });

export const fetchProcesses = (workspace) => {
  processes.update((req) => {
    switch (req.stage) {
      case 'idle':
        refetchProcesses(workspace);
        return pendingRequest();
      case 'pending':
        // Do not refetch.
        return req;
      case 'error':
        refetchProcesses(workspace);
        return pendingRequest();
      case 'success':
        refetchProcesses(workspace);
        return refetchingResponse(req.data);
    }
  });
}

export const startProcess = async (workspace, id: string) =>
  workspace.startProcess(id)
    .then(() => fetchProcesses(workspace));

export const stopProcess = async (workspace, id: string) =>
  workspace.stopProcess(id)
    .then(() => fetchProcesses(workspace));

export const refreshAllProcesses = async (workspace) =>
  workspace.refreshAllProcesses()
    .then(() => fetchProcesses(workspace));
