import { get, writable } from 'svelte/store';
import type { PaginationParams, WorkspaceApi } from '../api';
import {
  pendingRequest,
  refetchingResponse,
  RequestLifecycle,
  successResponse,
} from '../api';
import type { LogEvent } from './types';
import { visibleLogsStore } from './visible-logs';

const maxEvents = 1000;

export type LogsStore = Record<string, WorkspaceState>; // Keyed by workspaceId.

export interface WorkspaceState {
  cursor: string | null;
  filterStr: string | null;
  events: RequestLifecycle<LogEvent[]>;
}

export const logsStore = writable<LogsStore>({});

export const fetchLogs = async (
  workspaceId: string,
  workspace: WorkspaceApi,
  pagination: Partial<PaginationParams>,
  filterStr: string | null,
) => {
  const visibleLogs = [...get(visibleLogsStore).values()];
  if (visibleLogs.length == 0) {
    return;
  }

  logsStore.update((state) => {
    let workspaceState = state[workspaceId];
    if (workspaceState == null) {
      workspaceState = {
        cursor: null,
        filterStr: null,
        events: pendingRequest(),
      };
    }
    switch (workspaceState.events.stage) {
      case 'idle':
        workspaceState = {
          ...workspaceState,
          events: pendingRequest(),
        };
        break;
      case 'pending':
        // TODO: Prevent re-fetch of pending request.
        break;
      case 'error':
        workspaceState = {
          ...workspaceState,
          events: pendingRequest(),
        };
        break;
      case 'success':
        workspaceState = {
          ...workspaceState,
          events: refetchingResponse(workspaceState.events.data),
        };
        break;
    }
    return {
      ...state,
      [workspaceId]: workspaceState,
    };
  });

  // NOTE [FILTER BY LOG]:
  // The visible logs parameter acts as a filter for log events that restricts the results
  // to events that came from a stream belonging to only certain processes. This should be
  // replaced with a more flexible filtering framework.
  const newEvents = await workspace.getEvents(visibleLogs, filterStr, {
    cursor: get(logsStore)[workspaceId].cursor,
    ...pagination,
  });

  logsStore.update((state) => {
    const workspaceState = state[workspaceId]!;
    let prevEvents: LogEvent[] = [];
    if (
      workspaceState.events.stage === 'success' ||
      workspaceState.events.stage === 'refetching'
    ) {
      prevEvents = workspaceState.events.data;
    }
    const allEvents = [...prevEvents, ...newEvents.items];

    return {
      ...state,
      [workspaceId]: {
        ...workspaceState,
        cursor: newEvents.nextCursor,
        events: successResponse(allEvents.slice(allEvents.length - maxEvents)),
      },
    };
  });
};

export const refreshLogs = (workspaceId: string, workspace: WorkspaceApi) => {
  const state = get(logsStore);
  const { filterStr } = state[workspaceId];
  return fetchLogs(workspaceId, workspace, { next: 100 }, filterStr);
};

export const resetLogs = (workspaceId: string) =>
  logsStore.update((state) => ({
    ...state,
    [workspaceId]: {
      ...state[workspaceId],
      cursor: null,
      filterStr: '',
      events: successResponse([]),
    },
  }));

export const loadInitialLogs = (
  workspaceId: string,
  workspace: WorkspaceApi,
) => {
  resetLogs(workspaceId);
  return fetchLogs(workspaceId, workspace, { prev: 100 }, '');
};

export const setFilterStr = (workspaceId: string, workspace: WorkspaceApi, filterStr: string) => {
  logsStore.update((state) => ({
    ...state,
    [workspaceId]: {
      ...state[workspaceId],
      cursor: null,
      filterStr,
      events: successResponse([]),
    },
  }));
  return fetchLogs(workspaceId, workspace, { prev: 100 }, filterStr);
}
