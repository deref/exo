<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const variables = workspace.describeEnvironment();
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <Panel title="Workspace Variables" backRoute={workspaceRoute}>
    {#await variables}
      <Spinner />
    {:then variables}
      {#if Object.keys(variables).length === 0}
        <div>Empty Environment</div>
      {:else}
        <EnvironmentTable {variables} />
      {/if}
    {/await}
  </Panel>
</Layout>
