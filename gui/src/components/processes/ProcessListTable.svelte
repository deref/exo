<script lang="ts">
  import Icon from '../Icon.svelte';
  import IconButton from '../IconButton.svelte';
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
  <div class="card" style={logStyleFromHash(name + ':out')}>
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
        <Icon glyph="Ellipsis" />
      </IconButton>
      <div class="dropdown">
        <span>{name}</span>
        <button
          on:click={() => {
            router.push(
              `/workspaces/${encodeURIComponent(
                workspaceId,
              )}/processes/${encodeURIComponent(id)}`,
            );
          }}
        >
          <Icon glyph="Details" />
          View details
        </button>
        <button
          on:click={() => {
            setProcLogs(id, $visibleLogsStore.has(id) ? false : true);
          }}
        >
          <Icon glyph="Logs" />
          Toggle logs visibility
        </button>
        <button
          on:click={() => {
            void deleteProcess(workspace, id);
            setProcLogs(id, false);
          }}
        >
          <Icon glyph="Delete" />
          Remove from <b>exo</b>
        </button>
      </div>
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

  .dropdown {
    display: none;
    position: absolute;
    right: 0;
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    border-radius: 5px;
    padding: 4px 0;
    margin: -6px;
    z-index: 2;
  }

  .dropdown > span {
    display: block;
    padding: 4px 12px;
    font-size: 0.8em;
    color: var(--grey-7-color);
  }

  .dropdown button {
    background: none;
    border: none;
    display: flex;
    align-items: center;
    font-size: 0.9em;
    gap: 4px;
    border-radius: 2px;
    padding: 6px 18px;
    width: 100%;
    white-space: nowrap;
    color: var(--grey-5-color);
  }

  .dropdown button b {
    font-weight: 500;
    color: currentColor;
    color: var(--grey-3-color);
  }

  .dropdown button :global(*) {
    fill: currentColor;
  }

  .dropdown button :global(svg) {
    height: 16px;
    margin-left: -8px;
  }

  .dropdown button:hover,
  .dropdown button:hover b {
    color: var(--strong-color);
    background: var(--grey-e-color);
  }

  .actions {
    position: relative;
  }

  .actions:focus .dropdown,
  .actions:focus-within .dropdown {
    display: block;
  }
</style>
