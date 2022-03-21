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

  // load={workspace.describeVolumes}
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Storage" slot="navbar" />
  <Panel title="Volumes" backUrl={workspaceRoute}>
    <DataGrid
      columns={[
        {
          title: 'id',
          label: StringLabel,
          getValue: (volume) => volume.id,
        },
        {
          title: 'name',
          label: StringLabel,
          getValue: (volume) => volume.name,
        },
      ]}
      actions={[
        {
          tooltip: 'Edit volume',
          glyph: 'Edit',
          execute: (volume) => {
            router.push(
              `#/workspaces/${encodeURIComponent(
                workspaceId,
              )}/components/${encodeURIComponent(volume.id)}/edit`,
            );
          },
        },
        {
          tooltip: 'Delete volume',
          glyph: 'Delete',
          execute: async (volume) => {
            await workspace.deleteComponent(volume.id);
            window.location.reload();
          },
        },
      ]}
    />
  </Panel>
</Layout>
