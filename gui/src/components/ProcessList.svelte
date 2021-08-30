<script lang="ts">
  import Panel from './Panel.svelte';
  import IconButton from './IconButton.svelte';
  import ProcfileChecker from './processes/ProcfileChecker.svelte';
  import ProcessListTable from './processes/ProcessListTable.svelte';
  import AddSVG from './mono/AddSVG.svelte';
  import { onDestroy, onMount } from 'svelte';
  import type { RequestLifecycle, WorkspaceApi } from '../lib/api';
  import type { ProcessDescription } from '../lib/process/types';
  import {
    fetchProcesses,
    processes,
    refreshAllProcesses,
  } from '../lib/process/store';
  import * as router from 'svelte-spa-router';
  import RemoteData from './RemoteData.svelte';
  import IfEnabled from './IfEnabled.svelte';

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
</script>

<Panel
  title={`Workspace ${workspaceId}`}
  backRoute="/"
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
