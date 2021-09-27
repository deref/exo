<script lang="ts">
  import Panel from './Panel.svelte';
  import IconButton from './IconButton.svelte';
  import ProcfileChecker from './processes/ProcfileChecker.svelte';
  import ProcessListTable from './processes/ProcessListTable.svelte';
  import AddSVG from './mono/AddSVG.svelte';
  import DetailsSVG from './mono/DetailsSVG.svelte';
  import EllipsisSVG from './mono/EllipsisSVG.svelte';
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
        <IconButton tooltip="Workspace actions..." on:click={() => {}}>
          <EllipsisSVG />
        </IconButton>

        <div class="dropdown">
          <span>{displayName}</span>
          <button
            on:click={() => {
              router.push(
                `/workspaces/${encodeURIComponent(workspaceId)}/details`,
              );
            }}
          >
            <DetailsSVG />
            View details
          </button>
          <button
            on:click={() => {
              router.push(
                `#/workspaces/${encodeURIComponent(workspaceId)}/new-component`,
              );
            }}
          >
            <AddSVG />
            Add component
          </button>
        </div>
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
        <AddSVG /> Add component
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

  .dropdown button :global(*) {
    fill: currentColor;
  }

  .dropdown button :global(svg) {
    height: 16px;
    margin-left: -8px;
  }

  .dropdown button:hover {
    color: var(--strong-color);
    background: var(--grey-e-color);
  }

  .menu {
    position: relative;
  }

  .menu:focus .dropdown,
  .menu:focus-within .dropdown {
    display: block;
  }
</style>
