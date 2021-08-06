<script lang="ts">
  import { logStyleFromHash } from '../lib/color';
  import { onMount, onDestroy, afterUpdate, beforeUpdate } from 'svelte';
  import { hasData, isUnresolved } from '../lib/api';
  import type { WorkspaceApi } from '../lib/api';
  import { loadInitialLogs, logsStore, refreshLogs } from '../lib/logs/store';
  import type { WorkspaceState } from '../lib/logs/store';
  import { shortDate } from '../lib/time';
  import { processes } from '../lib/process/store';
  import FormattedLogMessage from './logs/FormattedLogMessage.svelte';

  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  const logsPollInterval = 1000;

  let state: WorkspaceState = {
    cursor: null,
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
</script>

<section>
  <h1>Logs</h1>
  {#if hasData(state.events)}
    <div class="log-overflow-wrapper">
      <div class="log-container" bind:this={logViewport}>
        {#each state.events.data as event (event.id)}
          <div class="log-entry" style={logStyleFromHash(event.log)}>
            <div class="timestamp">{shortDate(event.timestamp)}</div>
            <div class="log">{friendlyName(event.log)}</div>
            <div class="message">
              <FormattedLogMessage message={event.message} />
            </div>
          </div>
        {/each}
      </div>
    </div>
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

  .log-overflow-wrapper {
    overflow: hidden;
    border-radius: 4px;
    box-shadow: 0px 12px 16px -8px #00000033, 0px 0.25px 0px 1px #00000033;
  }

  .log-container {
    width: 100%;
    height: 100%;
    font-family: 'Fira Code', monospace;
    font-weight: 450;
    font-size: 15px;
    overflow-y: scroll;
  }

  .log-entry {
    display: flex;
  }

  .log-entry > div {
    padding: 0 0.3em;
    color: #333333;
  }

  .log-entry:hover {
    background: #f3f3f3;
    color: #111111;
  }

  .timestamp {
    color: #999999;
  }

  .log-entry:hover .timestamp {
    background: #eeeeee;
    color: #555555;
  }

  .log {
    background: var(--log-bg-color);
    color: var(--log-color);
  }

  .log-entry:hover .log {
    background: var(--log-bg-hover-color);
  }
</style>
