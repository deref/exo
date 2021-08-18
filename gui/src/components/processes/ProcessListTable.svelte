<script lang="ts">
  import IconButton from '../IconButton.svelte';
  import CheckboxButton from '../CheckboxButton.svelte';
  import RunSVG from '../mono/play.svelte';
  import StopSVG from '../mono/stop.svelte';
  import DeleteSVG from '../mono/delete.svelte';
  import { link } from 'svelte-spa-router';
  import {
    startProcess,
    stopProcess,
    deleteProcess,
  } from '../../lib/process/store';
  import {
    setLogVisibility,
    visibleLogsStore,
  } from '../../lib/logs/visible-logs';
  import type { ProcessDescription } from '../../lib/process/types';
  import type { WorkspaceApi } from '../../lib/api';

  export let data: ProcessDescription[];
  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  let statusPending = new Set<string>();

  function setProcLogs(processId: string, visible: boolean) {
    setLogVisibility(processId, visible);
  }

  function setProcRun(id: string, run: boolean) {
    statusPending = statusPending.add(id);
    const proc = data.find((p) => p.id === id);
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
</script>

<table>
  <thead>
    <th />
    <th>Process</th>
    <th>Logs</th>
    <th />
  </thead>
  {#each data as { id, name, running } (id)}
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
          tooltip={$visibleLogsStore.has(id) ? 'Hide logs' : 'Show logs'}
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

<style>
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
    padding: 4px 8px;
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
