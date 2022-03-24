<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import StringLabel from '../components/StringLabel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import DataGrid from '../components/DataGrid.svelte';
  import * as router from 'svelte-spa-router';
  import { query } from '../lib/graphql';

  export let params = { workspace: '' };
  const workspaceId = params.workspace;

  const q = query(
    `#graphql
    query ($workspaceId: String!) {
      workspace: workspaceById(id: $workspaceId) {
        id
        stack {
          networks {
            type
            name
            componentId
          }
        }
      }
    }`,
    {
      variables: {
        workspaceId,
      },
    },
  );

  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  $: networks = $q.data?.workspace?.stack?.networks;
</script>

<Layout loader={$q}>
  <WorkspaceNav {workspaceId} active="Networking" slot="navbar" />
  <Panel title="Networking" backUrl={workspaceRoute}>
    {#if networks}
      <DataGrid
        items={networks}
        plural="networks"
        columns={[
          {
            title: 'Component ID',
            label: StringLabel,
            getValue: (network) => network.componentId,
          },
          {
            title: 'Name',
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
                )}/components/${encodeURIComponent(network.componentId)}/edit`,
              );
            },
          },
          {
            tooltip: 'Delete network',
            glyph: 'Delete',
            execute: async (network) => {
              await workspace.deleteComponent(network.componentId);
              window.location.reload();
            },
          },
        ]}
      />
    {/if}
  </Panel>
</Layout>
