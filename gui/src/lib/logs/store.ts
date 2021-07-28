import { get, writable } from 'svelte/store';
import type { PaginationParams } from '../api';
import {
  notRequested,
  pendingRequest,
  refetchingResponse,
  RemoteData,
  successResponse,
} from '../api';
import type { LogEvent } from './types';
import { visibleLogsStore } from './visible-logs';

const maxEvents = 1000;

export type LogsStore = Record<string, WorkspaceState>; // Keyed by workspaceId.

export interface WorkspaceState {
  cursor: string | null;
  events: RemoteData<LogEvent[]>;
}

export const logsStore = writable<LogsStore>({});

export const fetchLogs = async (
  workspaceId: string,
  workspace,
  pagination: Partial<PaginationParams>,
) => {
  logsStore.update((state) => {
    let workspaceState = state[workspaceId];
    if (workspaceState == null) {
      workspaceState = {
        cursor: null,
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
  const visibleLogs = [...get(visibleLogsStore).values()];
  const newEvents = await workspace.getEvents(visibleLogs, {
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

export const refreshLogs = (workspaceId: string, workspace) =>
  fetchLogs(workspaceId, workspace, { next: 100 });

export const loadInitialLogs = (workspaceId, workspace) =>
  fetchLogs(workspaceId, workspace, { prev: 100 });

export const resetLogs = (workspaceId: string) => {
  logsStore.update((state) => {
    return {
      ...state,
      [workspaceId]: {
        ...state[workspaceId],
        cursor: null,
        events: successResponse([]),
      },
    };
  });
};
