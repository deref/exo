<script lang="ts">
  import Panel from './Panel.svelte';
  import IconButton from './IconButton.svelte';
  import CheckboxButton from './CheckboxButton.svelte';
  import AddSVG from './mono/add.svelte';
  import RunSVG from './mono/play.svelte';
  import StopSVG from './mono/stop.svelte';
  import DeleteSVG from './mono/delete.svelte';
  import { onDestroy, onMount } from 'svelte';
  import { link } from 'svelte-spa-router';
  import type { RemoteData, WorkspaceApi } from '../lib/api';
  import type { ProcessDescription } from '../lib/process/types';
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

    procfileExport = current === computed ? null : computed;
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

<Panel
  title={`Workspace ${workspaceId}`}
  backRoute="/"
  --panel-bg-color="var(--secondary-bg-color)"
  --panel-padding="0 1rem"
>
  <div slot="actions">
    <IconButton
      tooltip="Add new process"
      on:click={() => {
        router.push(
          `#/workspaces/${encodeURIComponent(workspaceId)}/new-process`,
        );
      }}
    >
      <AddSVG />
    </IconButton>
  </div>
  <section>
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
                <div class="run-controls">
                  {#if statusPending.has(id)}
                    <div class="spinner" />
                  {:else if running}
                    <div class="spinner running" />
                    <div class="control hover-only">
                      <IconButton
                        tooltip="Stop process"
                        on:click={() => setProcRun(id, false)}
                      >
                        <StopSVG />
                      </IconButton>
                    </div>
                  {:else}
                    <div class="stopped unhover-only" />
                    <div class="control hover-only">
                      <IconButton
                        tooltip="Run process"
                        on:click={() => setProcRun(id, true)}
                      >
                        <RunSVG />
                      </IconButton>
                    </div>
                  {/if}
                </div>
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
                  tooltip={$visibleLogsStore.has(id)
                    ? 'Hide logs'
                    : 'Show logs'}
                  on:click={() => {
                    setProcLogs(id, $visibleLogsStore.has(id) ? false : true);
                  }}
                  active={$visibleLogsStore.has(id)}
                />
              </td>

              <td>
                <div class="hover-only-visibility">
                  <IconButton
                    tooltip="Delete process"
                    on:click={() => {
                      void deleteProcess(workspace, id);
                      setProcLogs(id, false);
                    }}
                  >
                    <DeleteSVG />
                  </IconButton>
                </div>
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
          <button
            on:click={async () => {
              if (procfileExport == null) {
                return;
              }
              await workspace.writeFile('Procfile', procfileExport);
              checkProcfile();
            }}>Export?</button
          >
        </p>
      {/if}
    </div>
  </section>
</Panel>

<style>
  section {
    display: grid;
    grid-auto-flow: row;
    grid-template-rows: max-content 1fr;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    padding: 12px 0;
    color: var(--grey-7-color);
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

  .run-controls {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    margin-right: 8px;
  }

  tr:not(:hover):not(:focus-within) .hover-only {
    display: none;
  }

  tr:not(:hover):not(:focus-within) .hover-only-visibility {
    visibility: hidden;
  }

  tr:hover .unhover-only,
  tr:focus-within .unhover-only {
    display: none;
  }

  .stopped {
    width: 14px;
    height: 14px;
    border-radius: 2px;
    background: var(--grey-c-color);
  }

  .control {
    position: absolute;
    z-index: 4;
  }

  .spinner {
    position: absolute;
    z-index: 3;
    top: 7px;
    left: 7px;
    width: 18px;
    height: 18px;
    border-radius: 100%;
    animation: spin 1s infinite linear;
    border: 2px solid;
    border-top-color: var(--spinner-grey-t);
    border-right-color: var(--spinner-grey-r);
    border-bottom-color: var(--spinner-grey-b);
    border-left-color: var(--spinner-grey-l);
    transition: all 0.125s;
  }

  .spinner.running {
    border-top-color: var(--spinner-blue-t);
    border-right-color: var(--spinner-blue-r);
    border-bottom-color: var(--spinner-blue-b);
    border-left-color: var(--spinner-blue-l);
  }

  tr:hover .spinner.running,
  tr:focus-within .spinner.running {
    top: 2px;
    left: 2px;
    width: 28px;
    height: 28px;
  }

  .process-name {
    display: inline-block;
    text-decoration: none;
    margin: 0;
    line-height: 1;
    font-size: 16px;
    font-weight: 550;
    padding: 6px 9px;
    border-radius: 4px;
    color: var(--grey-5-color);
    background: var(--grey-d-color);
    outline: none;
  }

  .process-name:hover {
    color: var(--strong-color);
    background: var(--grey-c-color);
  }

  .process-name:focus {
    background: var(--grey-b-color);
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
