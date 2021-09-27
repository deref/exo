<script lang="ts">
  import IconButton from '../IconButton.svelte';
  import CheckboxButton from '../CheckboxButton.svelte';
  import ProcessRunControls from './ProcessRunControls.svelte';
  import DeleteSVG from '../mono/DeleteSVG.svelte';
  import EllipsisSVG from '../mono/EllipsisSVG.svelte';
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
  {#each data as { id, name, running } (id)}
    <tr>
      <td>
        <ProcessRunControls {setProcRun} {statusPending} {id} {running} />
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
        <div class="hover-half-visibility">
          <IconButton
            tooltip="Delete process"
            on:click={() => {
              void deleteProcess(workspace, id);
              setProcLogs(id, false);
            }}
          >
            <EllipsisSVG />
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

  td {
    font-size: inherit;
    font-weight: inherit;
    align-items: center;
    justify-content: center;
  }

  td:nth-child(2) {
    text-align: left;
  }

  td:not(:last-child):not(:first-child) {
    border-right: 16px solid transparent;
  }

  td:nth-child(2) {
    width: 99%;
  }

  td:not(:nth-child(2)) {
    white-space: nowrap;
  }

  table,
  td,
  tr {
    border: none;
  }

  tr:not(:last-child) {
    border-bottom: 8px solid transparent;
  }

  tr:not(:hover):not(:focus-within) .hover-half-visibility {
    opacity: 0.333;
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
</style>
