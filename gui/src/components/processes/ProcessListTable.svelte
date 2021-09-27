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

{#each data as { id, name, running } (id)}
  <div class="card">
    <div>
      <ProcessRunControls {setProcRun} {statusPending} {id} {running} />
    </div>

    <div>
      <a
        class="process-name"
        use:link
        href={`/workspaces/${encodeURIComponent(
          workspaceId,
        )}/processes/${encodeURIComponent(id)}`}
      >
        {name}
      </a>
    </div>

    <div class="checkbox">
      <CheckboxButton
        tooltip={$visibleLogsStore.has(id) ? 'Hide logs' : 'Show logs'}
        on:click={() => {
          setProcLogs(id, $visibleLogsStore.has(id) ? false : true);
        }}
        active={$visibleLogsStore.has(id)}
      />
    </div>

    <div class="actions">
      <IconButton>
        <EllipsisSVG />
      </IconButton>
      <div class="dropdown">
        <button
          on:click={() => {
            void deleteProcess(workspace, id);
            setProcLogs(id, false);
          }}
        >
          <DeleteSVG />
          Remove {name}
        </button>
      </div>
    </div>
  </div>
{:else}
  <i>No components yet.</i>
{/each}

<style>
  .card {
    box-shadow: var(--button-shadow);
    display: grid;
    grid-template-columns: max-content auto max-content max-content;
    align-items: center;
    padding: 4px;
    margin: 0px -4px;
    border-radius: 4px;
    margin-bottom: 8px;
  }

  .card:hover {
    box-shadow: var(--button-hover-shadow);
  }

  .card > * {
    font-size: inherit;
    font-weight: inherit;
    align-items: center;
    justify-content: center;
  }

  .card .checkbox {
    margin-right: 18px;
  }

  .card > *:nth-child(2) {
    text-align: left;
  }

  .card > *:not(:nth-child(2)) {
    white-space: nowrap;
  }

  .card:not(:hover):not(:focus-within) .actions {
    opacity: 0.333;
  }

  .actions {
    position: relative;
  }

  .actions button {
    background: none;
    border: none;
    display: flex;
    align-items: center;
    font-size: 0.9em;
    gap: 4px;
    border-radius: 2px;
    padding: 6px 18px;
    color: var(--grey-5-color);
  }

  .actions button :global(*) {
    fill: currentColor;
  }

  .actions button :global(svg) {
    height: 16px;
    margin-left: -8px;
  }

  .actions button:hover {
    color: var(--strong-color);
    background: var(--grey-e-color);
  }

  .actions .dropdown {
    display: none;
    position: absolute;
    right: 0;
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    border-radius: 5px;
    padding: 4px 0;
    z-index: 2;
  }

  .actions:focus .dropdown,
  .actions:focus-within .dropdown {
    display: block;
  }

  .process-name {
    display: inline-block;
    text-decoration: none;
    margin: 0;
    margin-left: 6px;
    margin-right: 12px;
    line-height: 1;
    font-size: 16px;
    font-weight: 550;
    padding: 6px 9px;
    border-radius: 3px;
    color: var(--grey-5-color);
    background: var(--grey-e-color);
    outline: none;
  }

  .process-name:hover {
    color: var(--strong-color);
    background: var(--grey-d-color);
  }

  .process-name:focus {
    background: var(--grey-d-color);
  }
</style>
