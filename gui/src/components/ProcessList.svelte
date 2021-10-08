<script lang="ts">
  import Icon from './Icon.svelte';
  import Panel from './Panel.svelte';
  import IfEnabled from './IfEnabled.svelte';
  import IconButton from './IconButton.svelte';
  import RemoteData from './RemoteData.svelte';
  import ContextMenu from './ContextMenu.svelte';
  import MenuItem from './MenuItem.svelte';
  import ProcfileChecker from './processes/ProcfileChecker.svelte';
  import ProcessListTable from './processes/ProcessListTable.svelte';
  import { onDestroy, onMount } from 'svelte';
  import type { RequestLifecycle, WorkspaceApi } from '../lib/api';
  import type { ProcessDescription } from '../lib/process/types';
  import {
    fetchProcesses,
    processes,
    refreshAllProcesses,
  } from '../lib/process/store';
  import * as router from 'svelte-spa-router';
  import Button from './Button.svelte';

  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  let processList: RequestLifecycle<ProcessDescription[]> = {
    stage: 'pending',
  };
  const unsubscribeProcesses = processes.subscribe((processes) => {
    processList = processes;
  });

  let refreshInterval: any;

  let procfileExport: string | null = null;
  async function checkProcfile() {
    const current = await workspace.readFile('Procfile');
    const computed = await workspace.exportProcfile();

    procfileExport = current && current === computed ? null : computed;
  }

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

  let modalOpen = false;

  const handleKeyDown = (ev: KeyboardEvent) => {
    if (ev.key === 'Escape') {
      modalOpen = false;
    }
  };

  const handleWrapClick = (_ev: MouseEvent) => {
    modalOpen = false;
  };
</script>

<svelte:window on:keydown={handleKeyDown} />

{#await workspace.describeSelf()}
  <Panel title="" backRoute="/" />
{:then description}
  <Panel title={description.displayName} backRoute="/" --panel-padding="0 1rem">
    <div class="actions" slot="actions">
      <span>Logs</span>
      <div class="menu" tabindex="0">
        <IconButton
          glyph="Ellipsis"
          tooltip="Workspace actions..."
          on:click={() => {}}
        />

        <ContextMenu title={description.displayName}>
          <MenuItem
            glyph="Details"
            href={`/workspaces/${encodeURIComponent(workspaceId)}/details`}
          >
            View details
          </MenuItem>
          <MenuItem
            glyph="Add"
            href={`#/workspaces/${encodeURIComponent(
              workspaceId,
            )}/new-component`}
          >
            Add component
          </MenuItem>
          <MenuItem
            glyph="Delete"
            danger
            on:click={() => {
              modalOpen = true;
            }}
          >
            Destroy workspace
          </MenuItem>
        </ContextMenu>
      </div>
    </div>

    {#if modalOpen}
      <div class="modal-wrap" on:click={handleWrapClick}>
        <div class="confirm" on:click|stopPropagation={() => {}}>
          <h3>Delete workspace?</h3>
          <p>
            Are you sure you want to delete {description.displayName}?<br />
            This is irreversible, but will only delete the workspace in exo, not
            the files.
          </p>
          <div class="buttons">
            <Button
              on:click={() => {
                modalOpen = false;
              }}>Cancel</Button
            >
            <Button
              danger
              on:click={() => {
                workspace.destroy();
                router.push('/');
                modalOpen = false;
              }}>Yes, delete</Button
            >
          </div>
        </div>
      </div>
    {/if}

    <section>
      <button
        id="add-component"
        on:click={() => {
          router.push(
            `#/workspaces/${encodeURIComponent(workspaceId)}/new-component`,
          );
        }}
      >
        <Icon glyph="Add" /> Add component
      </button>
      <RemoteData data={processList} let:data let:error>
        <div slot="success">
          <ProcessListTable {data} {workspace} {workspaceId} />
        </div>

        <div slot="error">
          Error fetching process list: {error}
        </div>
      </RemoteData>
      <IfEnabled feature="export procfile">
        <ProcfileChecker
          {procfileExport}
          clickHandler={async () => {
            if (procfileExport == null) {
              return;
            }
            await workspace.writeFile('Procfile', procfileExport);
            checkProcfile();
          }}
        />
      </IfEnabled>
    </section>
  </Panel>
{/await}

<style>
  .modal-wrap {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.25);
    z-index: 100;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .confirm {
    position: absolute;
    width: 480px;
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    border-radius: 6px;
    padding: 30px 36px;
  }

  .confirm p {
    margin: 0;
    margin-bottom: 24px;
  }

  .buttons {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 12px;
  }

  #add-component {
    background: none;
    font-size: 0.9em;
    color: var(--grey-5-color);
    border: 1px dashed var(--grey-c-color);
    width: calc(100% + 8px);
    display: flex;
    align-items: center;
    border-radius: 4px;
    gap: 6px;
    height: 40px;
    margin: 12px -4px;
    padding: 0 8px;
  }

  #add-component:hover,
  #add-component:focus,
  #add-component:focus-within {
    background: var(--grey-e-color);
    color: var(--strong-color);
  }

  #add-component :global(svg) {
    height: 16px;
  }

  #add-component :global(*) {
    fill: currentColor;
  }
  .actions {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-right: 16px;
  }
  .actions span {
    color: var(--grey-7-color);
  }

  .menu {
    outline: none;
    position: relative;
  }

  .menu:focus :global(nav),
  .menu:focus-within :global(nav) {
    display: block;
  }
</style>
