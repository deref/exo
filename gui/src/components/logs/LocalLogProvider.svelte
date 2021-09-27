<script lang="ts">
  import type { LogEvent } from './types';
  import ErrorLabel from '../ErrorLabel.svelte';
  import type { WorkspaceApi } from '../../lib/api';
  import { formatLogs } from './util';
  import { onDestroy } from 'svelte';

  const maxEvents = 1000;
  const logsPollInterval = 1000;

  let cursor: string | null = null;
  let events: LogEvent[];

  export let workspace: WorkspaceApi;

  export let filterStr: string | null = null;
  export let logs: string[] = [];
  export let processIdToName: Record<string, string> = {};

  let pollRefreshTimer: ReturnType<typeof setTimeout> | null = null;
  const scheduleNextPoll = async () => {
    // Only allow one refresh loop.
    if (pollRefreshTimer !== null) {
      return;
    }

    if (logs.length === 0) {
      cursor = null;
      events = [];
      return;
    }

    const res = await workspace.getEvents(logs, filterStr, {
      cursor,
      next: maxEvents,
    });
    cursor = res.nextCursor;
    events = [...events, ...formatLogs(res.items, processIdToName)].slice(
      -maxEvents,
    );

    pollRefreshTimer = setTimeout(() => {
      pollRefreshTimer = null;
      scheduleNextPoll();
    }, logsPollInterval);
  };

  onDestroy(() => {
    if (pollRefreshTimer !== null) {
      clearTimeout(pollRefreshTimer);
    }
  });

  const resetLogs = async (logs: string[], filterStr: string | null) => {
    if (logs.length === 0) {
      cursor = null;
      events = [];
      return;
    }

    const res = await workspace.getEvents(logs, filterStr, {
      cursor: null,
      prev: maxEvents,
    });
    cursor = res.nextCursor;
    events = formatLogs(res.items, processIdToName);
    scheduleNextPoll();
  };

  // Reset log events entirely when filter or logs change.
  $: {
    resetLogs(logs, filterStr);
  }
</script>

<slot {events}>
  <ErrorLabel value="Please provide a logs component" />
</slot>
