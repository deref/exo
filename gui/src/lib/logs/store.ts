import { get, writable } from 'svelte/store';
import type { PaginationParams } from '../api';
import { notRequested, pendingRequest, refetchingResponse, RemoteData, successResponse } from '../api';
import type { LogEvent } from './types';
import { visibleLogsStore } from './visible-logs';
export interface LogsStore {
  events: RemoteData<LogEvent[]>;
  logBufferSize: number;
}

let lastCursor: string | null = null;

export const logsStore = writable<LogsStore>({
  events: notRequested(),
  logBufferSize: 1000,
});

export const fetchLogs = async (workspace, pagination: Partial<PaginationParams>) => {
  logsStore.update((value) => {
    switch (value.events.stage) {
      case 'idle':
        return {
          ...value,
          events: pendingRequest(),
        };
      case 'pending':
        // TODO: Prevent re-fetch of pending request.
        return value;
      case 'error':
        return {
          ...value,
          events: pendingRequest(),
        };
      case 'success':
        return {
          ...value,
          events: refetchingResponse(value.events.data),
        };
    }
  });

  // NOTE [FILTER BY LOG]:
  // The visible logs parameter acts as a filter for log events that restricts the results
  // to events that came from a stream belonging to only certain processes. This should be
  // replaced with a more flexible filtering framework.
  const visibleLogs = [...get(visibleLogsStore).values()];
  const newEvents = await workspace.getEvents(visibleLogs, {
    cursor: lastCursor,
    ...pagination,
  });

  lastCursor = newEvents.nextCursor;
  logsStore.update(value => {
    let prevEvents: LogEvent[] = [];
    if (value.events.stage === 'success' || value.events.stage === 'refetching') {
      prevEvents = value.events.data;
    }
    const allEvents = [...prevEvents, ...newEvents.items];

    return {
      ...value,
      events: successResponse(allEvents.slice(allEvents.length-value.logBufferSize)),
    }
  });
};

export const refreshLogs = (workspace) => fetchLogs(workspace, { next: 100 });

export const loadInitialLogs = (workspace) => fetchLogs(workspace, { prev: 100 });

export const resetLogs = () => {
  lastCursor = null;
  logsStore.update((value) => {
    return {
      ...value,
      events: successResponse([]),
    };
  })
}
