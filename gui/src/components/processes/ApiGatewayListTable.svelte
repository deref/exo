<script lang="ts">
  import IconButton from '../IconButton.svelte';
  import ContextMenu from '../ContextMenu.svelte';
  import MenuItem from '../MenuItem.svelte';
  import ProcessRunControls from './ProcessRunControls.svelte';
  import * as router from 'svelte-spa-router';
  import {
    startProcess,
    stopProcess,
    deleteProcess,
  } from '../../lib/process/store';
  import { logStyleFromHash } from '../../lib/color';
  import type { WorkspaceApi, ApiGatewayDescription } from '../../lib/api';

  const { link } = router;

  export let data: ApiGatewayDescription[];
  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  let statusPending = new Set<string>();

  const setProcRun = async (id: string, run: boolean) => {
    statusPending = statusPending.add(id);
    const proc = data.find((p) => p.id === id);
    if (!proc) {
      console.error(`Cannot find process: ${id}`);
      return;
    }
    if (run) {
      await startProcess(workspace, id);
    } else {
      await stopProcess(workspace, id);
    }
    statusPending.delete(id);
  };
</script>

{#each data as { id, name, running } (id)}
  <div class="card" style={logStyleFromHash(name)}>
    <div>
      <ProcessRunControls
        {setProcRun}
        statusPending={statusPending.has(id)}
        {id}
        {running}
      />
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

    <div class="actions" tabindex="0">
      <IconButton glyph="Ellipsis" />

      <ContextMenu title={name}>
        <MenuItem
          glyph="Delete"
          danger
          on:click={() => {
            void deleteProcess(workspace, id);
          }}
        >
          Delete component
        </MenuItem>
      </ContextMenu>
    </div>
  </div>
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
