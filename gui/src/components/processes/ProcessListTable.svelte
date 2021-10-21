<script lang="ts">
  import IconButton from '../IconButton.svelte';
  import ContextMenu from '../ContextMenu.svelte';
  import MenuItem from '../MenuItem.svelte';
  import CheckboxButton from '../CheckboxButton.svelte';
  import ProcessRunControls from './ProcessRunControls.svelte';
  import * as router from 'svelte-spa-router';
  import {
    startProcess,
    stopProcess,
    deleteProcess,
  } from '../../lib/process/store';
  import {
    setLogVisibility,
    visibleLogsStore,
  } from '../../lib/logs/visible-logs';
  import { logStyleFromHash } from '../../lib/color';
  import type { ProcessDescription } from '../../lib/process/types';
  import type { WorkspaceApi } from '../../lib/api';

  const { link } = router;

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
  <div class="card" style={logStyleFromHash(name)}>
    <div>
      <ProcessRunControls {setProcRun} {statusPending} {id} {running} />
    </div>

    <div>
      <a
        class="process-name"
        use:link
        href={`/workspaces/${encodeURIComponent(
          workspaceId,
        )}/components/${encodeURIComponent(id)}`}
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

    <div class="actions" tabindex="0">
      <IconButton glyph="Ellipsis" />

      <ContextMenu title={name}>
        <MenuItem
          glyph="Details"
          href={`/workspaces/${encodeURIComponent(
            workspaceId,
          )}/components/${encodeURIComponent(id)}`}>View details</MenuItem
        >
        <MenuItem
          glyph="Edit"
          href={`/workspaces/${encodeURIComponent(
            workspaceId,
          )}/components/${encodeURIComponent(id)}/edit`}
          >Edit component</MenuItem
        >
        <MenuItem
          glyph="Logs"
          on:click={() => {
            setProcLogs(id, $visibleLogsStore.has(id) ? false : true);
          }}
        >
          Toggle logs visibility
        </MenuItem>
        <MenuItem
          glyph="Delete"
          danger
          on:click={() => {
            void deleteProcess(workspace, id);
            setProcLogs(id, false);
          }}
        >
          Delete component
        </MenuItem>
      </ContextMenu>
    </div>
  </div>
{:else}
  <i>No components yet.</i>
{/each}

<style>
  .card {
    box-shadow: var(--card-shadow);
    display: grid;
    grid-template-columns: max-content auto max-content max-content;
    align-items: center;
    padding: 4px;
    margin: 0px -4px;
    border-radius: 4px;
    margin-bottom: 8px;
    border-left: 2px solid var(--log-color);
  }

  .card:hover {
    box-shadow: var(--card-hover-shadow);
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

  .process-name {
    display: inline-block;
    text-decoration: none;
    margin: 0;
    margin-left: 6px;
    margin-right: 12px;
    line-height: 1;
    font-size: 16px;
    font-weight: 550;
    padding: 4px 7px;
    border-radius: 3px;
    color: var(--log-color);
    outline: none;
  }

  .process-name:hover,
  .process-name:focus {
    color: var(--log-hover-color);
    background: var(--log-bg-hover-color);
  }

  .actions {
    outline: none;
    position: relative;
  }

  .actions:focus :global(nav),
  .actions:focus-within :global(nav) {
    display: block;
  }
</style>
