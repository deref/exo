<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <Panel title="Workspace details" backRoute={workspaceRoute}>
    {#await workspace.describeSelf()}
      <Spinner />
    {:then description}
      <CheckeredTableWrapper>
        <table>
          <tbody>
            <tr>
              <td class="label">ID</td>
              <td>{description.id}</td>
            </tr>
            <tr>
              <td class="label">Display name</td>
              <td>{description.displayName}</td>
            </tr>
            <tr>
              <td class="label">Root directory</td>
              <td>{description.root}</td>
            </tr>
          </tbody>
        </table>
      </CheckeredTableWrapper>
    {/await}

    <!-- Databases, Apps, cloud services, etc. -->
  </Panel>
</Layout>

<style>
  .label {
    font-size: 0.8em;
    font-weight: 450;
    color: var(--grey-5-color);
  }
</style>
