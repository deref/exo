<script lang="ts">
  import Icon from './Icon.svelte';
  import Panel from './Panel.svelte';
  import IconButton from './IconButton.svelte';
  import ContextMenu from './ContextMenu.svelte';
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
  import { api } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import RemoteData from './RemoteData.svelte';
  import IfEnabled from './IfEnabled.svelte';

  export let workspace: WorkspaceApi;
  export let workspaceId: string;

  const workspaceName = api.kernel
    .describeWorkspaces()
    .then(
      (workspaces) =>
        workspaces.find((ws) => ws.id === workspaceId)?.displayName,
    );

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
</script>

{#await workspaceName}
  <Panel title="" backRoute="/" />
{:then displayName}
  <Panel title={displayName} backRoute="/" --panel-padding="0 1rem">
    <div class="actions" slot="actions">
      <span>Logs</span>
      <div class="menu">
        <IconButton
          glyph="Ellipsis"
          tooltip="Workspace actions..."
          on:click={() => {}}
        />
        <ContextMenu
          title={displayName}
          actions={[
            {
              name: 'View details',
              glyph: 'Details',
              function: () => {
                router.push(
                  `/workspaces/${encodeURIComponent(workspaceId)}/details`,
                );
              },
            },
            {
              name: 'Add component',
              glyph: 'Add',
              function: () => {
                router.push(
                  `#/workspaces/${encodeURIComponent(
                    workspaceId,
                  )}/new-component`,
                );
              },
            },
          ]}
        />
      </div>
    </div>
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
    position: relative;
  }

  .menu:focus :global(nav),
  .menu:focus-within :global(nav) {
    display: block;
  }
</style>
