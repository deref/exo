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
          stores {
            type
            name
            componentId
            sizeMiB
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

  $: stores = $q.data?.workspace?.stack?.stores;
</script>

<Layout loader={$q}>
  <WorkspaceNav {workspaceId} active="Storage" slot="navbar" />
  <Panel title="Storage" backUrl={workspaceRoute}>
    {#if stores}
      <DataGrid
        items={stores}
        plural="stores"
        columns={[
          {
            title: 'Component ID',
            label: StringLabel,
            getValue: (store) => store.componentId,
          },
          {
            title: 'Name',
            label: StringLabel,
            getValue: (store) => store.name,
          },
        ]}
        actions={[
          {
            tooltip: 'Edit store',
            glyph: 'Edit',
            execute: (store) => {
              router.push(
                `#/workspaces/${encodeURIComponent(
                  workspaceId,
                )}/components/${encodeURIComponent(store.componentId)}/edit`,
              );
            },
          },
          {
            tooltip: 'Delete store',
            glyph: 'Delete',
            execute: async (store) => {
              await workspace.deleteComponent(store.componentId);
              window.location.reload();
            },
          },
        ]}
      />
    {/if}
  </Panel>
</Layout>
