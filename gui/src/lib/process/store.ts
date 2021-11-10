import { writable } from 'svelte/store';
import {
  errorResponse,
  notRequested,
  pendingRequest,
  refetchingResponse,
  successResponse,
} from '../api';
import type { WorkspaceApi } from '../api';
import type { ProcessDescription } from './types';

export const processes = writable(notRequested<ProcessDescription[]>());

const refetchProcesses = (workspace: WorkspaceApi) =>
  workspace
    .describeProcesses()
    .then((data) => {
      processes.set(successResponse(data));
    })
    .catch((err: Error) => {
      processes.set(errorResponse(err.message));
    });

export const fetchProcesses = (workspace: WorkspaceApi) => {
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
      case 'refetching':
        // I'm not sure this is right.
        return req;
    }
  });
};

export const startProcess = async (workspace: WorkspaceApi, id: string) => {
  await workspace.startProcess(id);
  fetchProcesses(workspace);
};

export const stopProcess = async (workspace: WorkspaceApi, id: string) => {
  await workspace.stopProcess(id);
  fetchProcesses(workspace);
};

export const deleteProcess = async (workspace: WorkspaceApi, id: string) => {
  await workspace.deleteComponent(id);
  fetchProcesses(workspace);
};

export const refreshAllProcesses = async (workspace: WorkspaceApi) => {
  await workspace.refreshAllProcesses();
  fetchProcesses(workspace);
};
