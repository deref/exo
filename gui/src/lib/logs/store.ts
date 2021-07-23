import { get, writable } from 'svelte/store';
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

export const refreshLogs = async (workspace, fromStart = false) => {
  if (fromStart) {
    lastCursor = null;
  }

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
          events: refetchingResponse(fromStart ? [] : value.events.data),
        };
    }
  });

  const newEvents = await workspace.getEvents([...get(visibleLogsStore).values()], {
    cursor: fromStart ? null : lastCursor,
    next: 100,
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
