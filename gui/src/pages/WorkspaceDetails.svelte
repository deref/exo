<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api } from '../lib/api';
  import Spinner from '../components/Spinner.svelte';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <CenterFormPanel title="Workspace details" backRoute={workspaceRoute}>
    {#await workspace.describeSelf()}
      <Spinner />
    {:then description}
      <table>
        <tr>
          <th>ID</th>
          <td>{description.id}</td>
        </tr>
        <tr>
          <th>Display name</th>
          <td>{description.displayName}</td>
        </tr>
        <tr>
          <th>Root directory</th>
          <td>{description.root}</td>
        </tr>
      </table>
    {/await}

    <!-- Databases, Apps, cloud services, etc. -->
  </CenterFormPanel>
</Layout>

<style>
  th {
    text-align: left;
  }

  td {
    padding-left: 2em;
    padding-top: 0.3em;
    padding-bottom: 0.3em;
  }
</style>
