<script lang="ts">
import { onDestroy, onMount } from 'svelte';
import type { RemoteData } from '../lib/api';
import { logsStore, setLogVisibility } from '../lib/logs/store';
import { fetchProcesses, processes, startProcess, stopProcess, refreshAllProcesses, deleteProcess } from '../lib/process/store';
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

let loggedProcesses: string[] = [];
const unsubscribeLoggedProcesses = logsStore.subscribe(({ logs }) => {
  loggedProcesses = logs;
});

let processList: RemoteData<ProcessDescription[]> = { stage: 'pending' };
const unsubscribeProcesses = processes.subscribe(processes => {
  processList = processes;
  // TODO: Set all processes to be logged on the initial load. Subsequently,
  // the settings for whether to log processes should be persisted to LocalStorage
  // and should not be automatically refreshed.
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

function toggleProcLogs(id: string) {
  setLogVisibility(workspace, id, !loggedProcesses.includes(id));
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
  unsubscribeLoggedProcesses();
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
        {#if statusPending.has(name)}
        <button disabled><Loading /></button>
        {:else if running}
        <IconButton on:click={() => toggleProc(id)} activeState><Stop /></IconButton>
        {:else}
        <IconButton on:click={() => toggleProc(id)}><Run /></IconButton>
        {/if}

        {#if loggedProcesses.includes(id)}
        <IconButton on:click={() => toggleProcLogs(id)} activeState><Hide /></IconButton>
        {:else}
        <IconButton on:click={() => toggleProcLogs(id)}><Show /></IconButton>
        {/if}
        <IconButton on:click={() => deleteProcess(workspace, id)}><Delete/></IconButton>
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
  color: #bb0000;
  background: #ff000022;
}

p {
  font-weight: bold;
}

</style>
