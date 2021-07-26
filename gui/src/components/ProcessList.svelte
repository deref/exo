<script lang="ts">
import { onDestroy, onMount } from 'svelte';
import type { RemoteData } from '../lib/api';
import { loadInitialLogs, resetLogs } from '../lib/logs/store';
import { fetchProcesses, processes, startProcess, stopProcess, refreshAllProcesses, deleteProcess } from '../lib/process/store';
import { toggleLogVisibility, visibleLogsStore } from '../lib/logs/visible-logs';
import type { ProcessDescription } from '../lib/process/types';
import * as router from 'svelte-spa-router';
import IconButton from './IconButton.svelte';

//Icons
import Add from './mono/add.svelte';
import Run from './mono/play.svelte'
import Loading from './mono/refresh.svelte'
import Stop from './mono/stop.svelte'
import Show from './mono/eye.svelte'
import Hide from './mono/eye-off.svelte'
import Delete from './mono/delete.svelte'

export let workspace;
export let workspaceId: string;

let statusPending = new Set<string>();

let processList: RemoteData<ProcessDescription[]> = { stage: 'pending' };
const unsubscribeProcesses = processes.subscribe(processes => {
  processList = processes;
});

function toggleProc(id: string) {
  if (processList.stage !== 'success') {
    return;
  }
  statusPending = statusPending.add(id);
  const proc = processList.data.find(p => p.id === id);
  if (!proc) {
    console.error(`Cannot find process: ${id}`);
    return;
  }
  if (proc.running) {
    stopProcess(workspace, id).then(() => {
      statusPending.delete(id);
    });
  } else {
    startProcess(workspace, id).then(() => {
      statusPending.delete(id);
    });
  }
}

function toggleProcLogs(processId: string) {
  toggleLogVisibility(processId);
  resetLogs();
  loadInitialLogs(workspace);
}

onMount(() => {
  fetchProcesses(workspace);
  
  // TODO: Server-sent events or websockets!
  setInterval(() => {
    refreshAllProcesses(workspace);
  }, 5000);
});

onDestroy(() => {
  unsubscribeProcesses();
});
</script>

<section>
  <h1>
    Processes
    <IconButton on:click={() => {
      router.push(`#/workspaces/${encodeURIComponent(workspaceId)}/new-process`)
    }}>
      <Add/>
    </IconButton>
  </h1>
  <div>
  {#if processList.stage == 'pending' || processList.stage == 'idle'}
    Loading...
  {:else if processList.stage == 'success' || processList.stage == 'refetching'}
    {#each processList.data as {id, name, running} (id)}
    <div class="process-description">
      <h2>{name}</h2>
      <div></div>
      <div class="actions">
        {#if statusPending.has(id)}
        <button disabled><Loading /></button>
        {:else if running}
        <IconButton tooltip="Stop process" on:click={() => toggleProc(id)} active><Stop /></IconButton>
        {:else}
        <IconButton tooltip="Run process" on:click={() => toggleProc(id)}><Run /></IconButton>
        {/if}

        {#if $visibleLogsStore.has(id)}
        <IconButton tooltip="Hide logs" on:click={() => toggleProcLogs(id)} active><Hide /></IconButton>
        {:else}
        <IconButton tooltip="Show logs" on:click={() => toggleProcLogs(id)}><Show /></IconButton>
        {/if}
        <IconButton tooltip="Delete process" on:click={() => deleteProcess(workspace, id)}><Delete/></IconButton>
      </div>
    </div>
    {:else}
    <i>No processes yet.</i>
    {/each}
  {:else if processList.stage == 'error'}
    Error fetching process list: {processList.message}
  {/if}
  </div>
</section>

<style>

section {
  display: grid;
  grid-auto-flow: row;
  grid-template-rows: max-content 1fr;
}

.process-description {
  display: grid;
  grid-template-columns: max-content auto max-content;
  gap: 12px;
  margin-bottom: 8px;
  align-items: center;
}

h1 {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

h2 {
  margin: 0;
  line-height: 1;
  font-size: 18px;
  font-weight: 550;
  padding: 6px 12px;
  border-radius: 6px;
  color: #555;
  background: #eee;
}

</style>
