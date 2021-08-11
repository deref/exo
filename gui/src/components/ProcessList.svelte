<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { link } from 'svelte-spa-router';
  import type { RemoteData, WorkspaceApi } from '../lib/api';
  import { loadInitialLogs, resetLogs } from '../lib/logs/store';
  import {
    fetchProcesses,
    processes,
    startProcess,
    stopProcess,
    refreshAllProcesses,
    deleteProcess,
  } from '../lib/process/store';
  import { setLogVisibility, visibleLogsStore } from '../lib/logs/visible-logs';
  import * as router from 'svelte-spa-router';
  import IconButton from './IconButton.svelte';

  //Icons
  import Add from './mono/add.svelte';
  import Run from './mono/play.svelte';
  import Loading from './mono/refresh.svelte';
  import Stop from './mono/stop.svelte';
  import Delete from './mono/delete.svelte';
  import CheckboxButton from './CheckboxButton.svelte';
  import type { ProcessDescription } from 'src/lib/process/types';

  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  let statusPending = new Set<string>();

  let processList: RemoteData<ProcessDescription[]> = { stage: 'pending' };
  const unsubscribeProcesses = processes.subscribe((processes) => {
    processList = processes;
  });

  function setProcRun(id: string, run: boolean) {
    if (processList.stage !== 'success') {
      return;
    }
    statusPending = statusPending.add(id);
    const proc = processList.data.find((p) => p.id === id);
    if (!proc) {
      console.error(`Cannot find process: ${id}`);
      return;
    }
    if (run) {
      startProcess(workspace, id).then(() => {
        statusPending.delete(id);
      });
    } else {
      stopProcess(workspace, id).then(() => {
        statusPending.delete(id);
      });
    }
  }

  function setProcLogs(processId: string, visible: boolean) {
    setLogVisibility(processId, visible);
    resetLogs(workspaceId);
    loadInitialLogs(workspaceId, workspace);
  }

  let procfileExport: string | null = null;
  async function checkProcfile() {
    const current = await workspace.readFile('Procfile');
    const computed = await workspace.exportProcfile();

    procfileExport = (current === computed) ? null : computed;
  }

  let refreshInterval: any;

  onMount(() => {
    fetchProcesses(workspace);
    checkProcfile();

    // TODO: Server-sent events or websockets!
    refreshInterval = setInterval(() => {
      refreshAllProcesses(workspace);
      checkProcfile();
    }, 5000);
  });

  onDestroy(() => {
    clearInterval(refreshInterval);
    unsubscribeProcesses();
  });
</script>

<section>
  <h1>
    Processes
    <IconButton
      tooltip="Add new process"
      on:click={() => {
        router.push(
          `#/workspaces/${encodeURIComponent(workspaceId)}/new-process`,
        );
      }}
    >
      <Add />
    </IconButton>
  </h1>
  <div>
    {#if processList.stage == 'pending' || processList.stage == 'idle'}
      Loading...
    {:else if processList.stage == 'success' || processList.stage == 'refetching'}
      <table>
        <thead>
          <th />
          <th>Process</th>
          <th>Logs</th>
          <th />
        </thead>
        {#each processList.data as { id, name, running } (id)}
          <tr>
            <td>
              {#if statusPending.has(id)}
                <button disabled><Loading /></button>
              {:else if running}
                <IconButton
                  tooltip="Stop process"
                  on:click={() => setProcRun(id, false)}
                  active><Stop /></IconButton
                >
              {:else}
                <IconButton
                  tooltip="Run process"
                  on:click={() => setProcRun(id, true)}><Run /></IconButton
                >
              {/if}
            </td>

            <td
              ><a
                class="process-name"
                use:link
                href={`/workspaces/${encodeURIComponent(
                  workspaceId,
                )}/processes/${encodeURIComponent(id)}`}>{name}</a
              ></td
            >

            <td>
              <CheckboxButton
                tooltip={$visibleLogsStore.has(id) ? 'Hide logs' : 'Show logs'}
                on:click={() => {
                  setProcLogs(id, $visibleLogsStore.has(id) ? false : true);
                }}
                active={$visibleLogsStore.has(id)}
              />
            </td>

            <td>
              <IconButton
                tooltip="Delete process"
                on:click={() => {
                  void deleteProcess(workspace, id);
                  setProcLogs(id, false);
                }}><Delete /></IconButton
              >
            </td>
          </tr>
        {:else}
          <i>No processes yet.</i>
        {/each}
      </table>
    {:else if processList.stage == 'error'}
      Error fetching process list: {processList.message}
    {/if}
  </div>
  <div>
    {#if procfileExport}
      <p>
        Your Procfile is not up to date.
        <button on:click={async () => {
          if (procfileExport == null) {
            return;
          }
          await workspace.writeFile('Procfile', procfileExport);
          checkProcfile();
        }}>Export?</button>
      </p>
    {/if}
  </div>
</section>

<style>
  section {
    display: grid;
    grid-auto-flow: row;
    grid-template-rows: max-content 1fr;
  }

  table {
    width: calc(100% + 12px);
    border-collapse: collapse;
    margin-left: -12px;
  }

  th {
    padding: 12px 0;
  }

  td,
  th {
    font-size: inherit;
    font-weight: inherit;
    align-items: center;
    justify-content: center;
  }

  td:nth-child(2),
  th:nth-child(2) {
    text-align: left;
  }

  td:not(:last-child):not(:first-child),
  th:not(:last-child):not(:first-child) {
    border-right: 16px solid transparent;
  }

  td:nth-child(2) {
    width: 99%;
  }

  td:not(:nth-child(2)) {
    white-space: nowrap;
  }

  table,
  thead,
  th,
  td,
  tr {
    border: none;
  }

  tr:not(:last-child) {
    border-bottom: 8px solid transparent;
  }

  h1 {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .process-name {
    display: inline-block;
    margin: 0;
    line-height: 1;
    font-size: 16px;
    font-weight: 550;
    padding: 8px 12px;
    border-radius: 4px;
    color: #555;
    background: #eee;
  }
</style>
