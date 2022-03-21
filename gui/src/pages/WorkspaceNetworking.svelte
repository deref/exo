<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import StringLabel from '../components/StringLabel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import DataGrid from '../components/DataGrid.svelte';
  import { api } from '../lib/api';
  import * as router from 'svelte-spa-router';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  // load={workspace.describeNetworks}
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Networking" slot="navbar" />
  <Panel title="Networks" backUrl={workspaceRoute}>
    <DataGrid
      columns={[
        {
          title: 'id',
          label: StringLabel,
          getValue: (network) => network.id,
        },
        {
          title: 'name',
          label: StringLabel,
          getValue: (network) => network.name,
        },
      ]}
      actions={[
        {
          tooltip: 'Edit network',
          glyph: 'Edit',
          execute: (network) => {
            router.push(
              `#/workspaces/${encodeURIComponent(
                workspaceId,
              )}/components/${encodeURIComponent(network.id)}/edit`,
            );
          },
        },
        {
          tooltip: 'Delete network',
          glyph: 'Delete',
          execute: async (network) => {
            await workspace.deleteComponent(network.id);
            window.location.reload();
          },
        },
      ]}
    />
  </Panel>
</Layout>
