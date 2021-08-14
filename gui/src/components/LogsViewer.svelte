<script lang="ts">
  import Panel from './Panel.svelte';
  import Textbox from './Textbox.svelte';
  import FormattedLogMessage from './logs/FormattedLogMessage.svelte';
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
    <div class="log-table-container" bind:this={logViewport}>
      <table>
        {#each state.events.data as event (event.id)}
          <tr class="log-entry" style={logStyleFromHash(event.log)}>
            <td class="timestamp">
              <span class="short-time">{shortTime(event.timestamp)}</span>
              <span class="full-timestamp">{event.timestamp}</span>
            </td>
            <td>{friendlyName(event.log)}</td>
            <td>
              <FormattedLogMessage message={event.message} />
            </td>
          </tr>
        {/each}
      </table>
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
  .log-table-container {
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

  table {
    font-family: var(--font-mono);
    font-variant-ligatures: var(--preferred-ligatures-logs);
    font-weight: 450;
    font-size: 15px;
  }

  table,
  tr,
  td {
    border: none;
    border-collapse: collapse;
  }

  td {
    padding: 0 0.3em;
    vertical-align: text-top;
    color: #333333;
    white-space: pre-wrap;
  }

  tr:hover td {
    background: #f3f3f3;
    color: #111111;
  }

  td:nth-child(1) {
    color: #999999;
  }

  tr:hover td:nth-child(1) {
    background: #eeeeee;
    color: #555555;
  }

  td:nth-child(2) {
    text-align: right;
    background: var(--log-bg-color);
    color: var(--log-color);
  }

  tr:hover td:nth-child(2) {
    background: var(--log-bg-hover-color);
    color: var(--log-hover-color);
  }

  .timestamp {
    position: relative;
  }

  .full-timestamp {
    display: none;
  }

  .timestamp:hover {
    overflow: visible;
  }

  .timestamp:hover .full-timestamp {
    display: block;
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: auto;
    white-space: nowrap;
    background: #333333;
    color: #f3f3f3;
    padding: 0 0.3em;
  }
</style>
