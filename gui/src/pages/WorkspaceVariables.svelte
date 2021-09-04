<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { api } from '../lib/api';
  import Spinner from '../Components/Spinner.svelte';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const variables = workspace.describeEnvironment();
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Components" slot="navbar" />
  <Panel title="Variables" backRoute={workspaceRoute}>
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
