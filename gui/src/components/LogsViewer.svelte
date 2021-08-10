<script lang="ts">
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
  import { shortDate } from '../lib/time';
  import { processes } from '../lib/process/store';
  import FormattedLogMessage from './logs/FormattedLogMessage.svelte';
  import Textbox from './Textbox.svelte';
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
</script>

<section>
  <h1>Logs</h1>
  {#if hasData(state.events)}
    <div class="log-table-overflow-wrapper">
      <div class="log-table-container" bind:this={logViewport}>
        <table>
          {#each state.events.data as event (event.id)}
            <tr class="log-entry" style={logStyleFromHash(event.log)}>
              <td>{shortDate(event.timestamp)}</td>
              <td>{friendlyName(event.log)}</td>
              <td>
                <FormattedLogMessage message={event.message} />
              </td>
            </tr>
          {/each}
        </table>
      </div>
    </div>
    <Textbox
      placeholder="Filter..."
      value={state.filterStr || ''}
      on:input={(e) => {
        const text = e.currentTarget?.value.trim();
        setFilterStrDebounced(text);
      }}
      --input-width="100%"
    />
  {:else if isUnresolved(state.events)}
    <div>Loading logs...</div>
  {:else}
    <div>Error fetching logs: {state.events.message}</div>
  {/if}
</section>

<style>
  section {
    overflow: hidden;
    padding: 1px;
    display: grid;
    grid-auto-flow: row;
    grid-template-rows: max-content 1fr;
  }

  .log-table-overflow-wrapper {
    overflow: hidden;
    border-radius: 4px;
    box-shadow: 0px 12px 16px -8px #00000033, 0px 0.25px 0px 1px #00000033;
    margin-bottom: 18px;
  }

  .log-table-container {
    width: 100%;
    height: 100%;
    overflow-x: auto;
    overflow-y: scroll;
  }

  table {
    font-family: 'Fira Code', monospace;
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
  }
</style>
