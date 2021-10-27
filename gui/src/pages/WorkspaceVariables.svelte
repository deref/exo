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

  const makeRequests = () =>
    Promise.all([
      workspace.describeEnvironment(),
      // workspace.describeVaults()
    ]);
  let requests = makeRequests();
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <Panel title="Workspace Variables" backRoute={workspaceRoute}>
    {#await requests}
      <Spinner />
    {:then [variables]}
      {#if Object.keys(variables).length > 0}
        <h2>Variables</h2>
        <EnvironmentTable
          variables={Object.entries(variables).map(([name, description]) => ({
            name,
            ...description,
          }))}
        />
      {:else}
        <div>Empty environment, no variables found.</div>
      {/if}
    {/await}
  </Panel>
</Layout>
