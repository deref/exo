<script lang="ts">
  
import { hashPalette } from '../lib/color';
import { onMount, onDestroy, afterUpdate, beforeUpdate } from 'svelte';
import type { RemoteData } from '../lib/api';
import { hasData, IsUnresolved, notRequested } from '../lib/api';
import { logsStore, refreshLogs } from '../lib/logs/store';
import type { LogEvent } from '../lib/logs/types';
import { shortDate } from '../lib/time';

export let workspace;

const logsPollInterval = 1000;

let logEvents: RemoteData<LogEvent[]> = notRequested();
const unsubscribeProcesses = logsStore.subscribe(data => {
  logEvents = data.events;
});

// Poll server for new logs.
let pollRefreshTimer = null;
const scheduleNextPoll = () => {
  refreshLogs(workspace);
  pollRefreshTimer = setTimeout(scheduleNextPoll, logsPollInterval);
};

onMount(() => {
  scheduleNextPoll();
});

onDestroy(() => {
  unsubscribeProcesses();
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
  const currentPosition = logViewport.scrollTop +logViewport.offsetHeight;
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
  {#if hasData(logEvents)}
  <div class="log-table-overflow-wrapper">
    <div class="log-table-container" bind:this={logViewport}>
      <table>
        {#each logEvents.data as event (event.id)}
        <tr class="log-entry" style={`--log-color:${hashPalette(event.log)[0]};--log-bg-color:${hashPalette(event.log)[2]};--log-border-color:${hashPalette(event.log)[1]}`}>
            <td>{event.log}</td> <td>{shortDate(event.timestamp)}</td> <td>{event.message}</td>
        </tr>
        {/each}
      </table>
    </div>
  </div>
  {:else if IsUnresolved(logEvents)}
    <div>Loading logs...</div>
  {:else}
    <div>Error fetching logs: {logEvents.message}</div>
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

table, tr, td {
  border: none;
  border-collapse: collapse;
}

td {
  padding: 0 0.4em;
  white-space: nowrap;
}

td:nth-child(1) {
  background: var(--log-bg-color);
  color: var(--log-color);
  border-right: 1px solid var(--log-border-color);
}

td:nth-child(2) {
  background: #77777711;
  color: #777777;
  border-right: 1px solid #77777777;
}
</style>
