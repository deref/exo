<script lang="ts">
import { onDestroy, onMount } from 'svelte';
import type { RemoteData } from '../lib/api';
import { logsStore, setLogVisibility } from '../lib/logs/store';
import { fetchProcesses, processes, startProcess, stopProcess, refreshAllProcesses } from '../lib/process/store';
import type { ProcessDescription } from '../lib/process/types';

//Icons
import Run from './mono/play.svelte'
import Loading from './mono/refresh.svelte'
import Stop from './mono/stop.svelte'
import Show from './mono/eye.svelte'
import Hide from './mono/eye-off.svelte'

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
    stopProcess(id).then(() => {
      statusPending.delete(id);
    });
  } else {
    startProcess(id).then(() => {
      statusPending.delete(id);
    });
  }
}

function toggleProcLogs(id: string) {
  setLogVisibility(id, !loggedProcesses.includes(id));
}

onMount(() => {
  fetchProcesses();
});

onDestroy(() => {
  unsubscribeProcesses();
  unsubscribeLoggedProcesses();
});
</script>

<section>
  <h1>Processes <a on:click={refreshAllProcesses}>(refresh)</a></h1>
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
        <button on:click={() => toggleProc(id)}><Stop /></button>
        {:else}
        <button on:click={() => toggleProc(id)}><Run /></button>
        {/if}

        {#if loggedProcesses.includes(id)}
        <button on:click={() => toggleProcLogs(id)}><Hide /></button>
        {:else}
        <button on:click={() => toggleProcLogs(id)}><Show /></button>
        {/if}
      </div>
    </div>
    {/each}
  {:else if processList.stage == 'error'}
    <p>Error fetching process list: {processList.message}</p>
  {/if}
</section>

<style>

  .process-description {
    display: grid;
    grid-template-columns: max-content auto max-content;
    gap: 12px;
    margin-bottom: 8px;
    align-items: center;
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

  h1 a {
    font-size: 14px;
    color: blue;
  }

  button {
    font-family: inherit;
    font-size: 0;
    font-weight: 450;
    padding: 6px;
    color: #777777;
    background-color: #77777700;
    border-radius: 4px;
    border: none;
    outline: none;
    margin-left: 4px;
  }

  button :global(svg), button :global(svg *) {
    fill: currentColor;
  }

  button:hover {
    background-color: #77777711;
    color: #444444;
  }

  button:focus {
    background-color: #77777733;
  }

  p {
    font-weight: bold;
  }
</style>
