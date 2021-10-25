<script lang="ts">
  import type { LogEvent } from '../../lib/logs/types';
  import ErrorLabel from '../ErrorLabel.svelte';
  import type { WorkspaceApi } from '../../lib/api';
  import { onDestroy } from 'svelte';

  const maxEvents = 1000;
  const pollInterval = 1000;

  let cursor: string | null = null;
  let events: LogEvent[];

  export let workspace: WorkspaceApi;

  export let filterStr: string | null = null;
  export let streams: string[] = [];

  let pollRefreshTimer: ReturnType<typeof setTimeout> | null = null;
  const scheduleNextPoll = async () => {
    // Only allow one refresh loop.
    if (pollRefreshTimer !== null) {
      return;
    }

    const res = await workspace.getEvents(streams, filterStr, {
      cursor,
      next: maxEvents,
    });
    cursor = res.nextCursor;
    events = [...events, ...res.items].slice(-maxEvents);

    pollRefreshTimer = setTimeout(() => {
      pollRefreshTimer = null;
      scheduleNextPoll();
    }, pollInterval);
  };

  onDestroy(() => {
    if (pollRefreshTimer !== null) {
      clearTimeout(pollRefreshTimer);
    }
  });

  const resetStreams = async (streams: string[], filterStr: string | null) => {
    const res = await workspace.getEvents(streams, filterStr, {
      cursor: null,
      prev: maxEvents,
    });
    cursor = res.nextCursor;
    events = res.items;
    scheduleNextPoll();
  };

  // Reset log events entirely when filter or streams change.
  $: {
    resetStreams(streams, filterStr);
  }

  // This is not ideal, since any change to the set of displayed logs will
  // reset the streams and the cursor, effectively "unclearing" events.
  // However, even this simple implementation should suffice for a while to
  // satisfy the use case of clearing scrollback to make it easier to see only
  // recent events while debugging.
  export const clearEvents = () => {
    events = [];
  };
</script>

<slot {events} {clearEvents}>
  <ErrorLabel value="Please provide a logs component" />
</slot>
