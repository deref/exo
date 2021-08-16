<script lang="ts">
  import Panel from './Panel.svelte';
  import Logs from './logs/Logs.svelte';
  import { logStyleFromHash } from '../lib/color';
  import { onMount, onDestroy, afterUpdate, beforeUpdate } from 'svelte';
  import { hasData, isUnresolved } from '../lib/api';
  import type { WorkspaceApi } from '../lib/api';
  import {
    loadInitialLogs,
    logsStore,
    refreshLogs,
    setFilterStr,
  } from '../lib/logs/store';
  import type { WorkspaceState } from '../lib/logs/store';
  import { shortTime } from '../lib/time';
  import { processes } from '../lib/process/store';
  import debounce from '../lib/debounce';

  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  const logsPollInterval = 1000;

  let state: WorkspaceState = {
    cursor: null,
    filterStr: null,
    events: { stage: 'success', data: [] },
  };
  const unsubscribeLogStore = logsStore.subscribe((workspaces) => {
    state = workspaces[workspaceId] ?? state;
  });

  // Poll server for new logs.
  let pollRefreshTimer: ReturnType<typeof setTimeout> | null = null;
  const scheduleNextPoll = async () => {
    await refreshLogs(workspaceId, workspace);
    pollRefreshTimer = setTimeout(scheduleNextPoll, logsPollInterval);
  };

  let knownProcessNameById: Record<string, string> = {};
  const unsubscribeProcessStore = processes.subscribe((processData) => {
    if (processData.stage === 'success' || processData.stage === 'refetching') {
      knownProcessNameById = processData.data.reduce(
        (acc, processDescription) => ({
          ...acc,
          [processDescription.id]: processDescription.name,
        }),
        {},
      );
    }
  });

  const friendlyName = (log: string): string => {
    // SEE NOTE [LOG_COMPONENTS].
    const [procId, stream] = log.split(':');
    const procName = knownProcessNameById[procId];
    return procName ? (stream ? `${procName}:${stream}` : procName) : log;
  };

  const formatLogs = (events: any) => {
    return events.data.map((event: any) => {
      return {
        id: event.id,
        style: logStyleFromHash(event.log),
        time: {
          short: shortTime(event.timestamp),
          full: event.timestamp,
        },
        name: friendlyName(event.log),
        message: event.message,
      };
    });
  };

  onMount(() => {
    loadInitialLogs(workspaceId, workspace).then(scheduleNextPoll);
  });

  onDestroy(() => {
    unsubscribeLogStore();
    unsubscribeProcessStore();
    if (pollRefreshTimer !== null) {
      clearTimeout(pollRefreshTimer);
    }
  });

  let logViewport: HTMLElement;
  let wasScrolledCloseToBottom = true;
  // Record whether the user was scrolled close to the bottom before new entries arrived.
  // If so, scroll them to the new bottom after the update.
  beforeUpdate(() => {
    if (!logViewport) {
      return;
    }

    const threshold = 150;
    const currentPosition = logViewport.scrollTop + logViewport.offsetHeight;
    const height = logViewport.scrollHeight;
    wasScrolledCloseToBottom = currentPosition > height - threshold;
  });

  afterUpdate(async () => {
    if (wasScrolledCloseToBottom && logViewport) {
      logViewport.scrollTop = logViewport.scrollHeight;
    }
  });

  const setFilterStrDebounced = debounce((filterStr: string) => {
    setFilterStr(workspaceId, workspace, filterStr);
  }, 250);
  let filterInput = state.filterStr;
  $: {
    if (filterInput !== null) {
      const text = filterInput.trim();
      setFilterStrDebounced(text);
    }
  }
</script>

<Panel title="Logs" --panel-padding="0" --panel-overflow-y="hidden">
  {#if hasData(state.events)}
    <div class="logs-container" bind:this={logViewport}>
      <Logs logs={formatLogs(state.events)} />
    </div>
  {:else if isUnresolved(state.events)}
    <div>Loading logs...</div>
  {:else}
    <div>Error fetching logs: {state.events.message}</div>
  {/if}
  <div slot="bottom">
    <input type="text" placeholder="Filter..." bind:value={filterInput} />
  </div>
</Panel>

<style>
  .logs-container {
    width: 100%;
    height: 100%;
    overflow-y: scroll;
    overflow-x: hidden;
  }

  input {
    width: 100%;
    background: var(--primary-bg-color);
    border: none;
    border-top: 1px solid var(--layout-bg-color);
    font-size: 16px;
    padding: 8px 12px;
    outline: none;
  }
</style>
