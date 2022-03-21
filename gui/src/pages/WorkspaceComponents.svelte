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
          components {
            id
            name
            type
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

  $: components = $q.data?.workspace?.stack?.components;
</script>

<Layout loading={$q.loading} error={$q.error}>
  <WorkspaceNav {workspaceId} active="Components" slot="navbar" />
  <Panel title="Components" backUrl={workspaceRoute}>
    {#if components}
      <DataGrid
        items={components}
        columns={[
          {
            title: 'id',
            label: StringLabel,
            getValue: (component) => component.id,
          },
          {
            title: 'name',
            label: StringLabel,
            getValue: (component) => component.name,
          },
          {
            title: 'type',
            label: StringLabel,
            getValue: (component) => component.type,
          },
        ]}
        actions={[
          {
            tooltip: 'Edit component',
            glyph: 'Edit',
            execute: (component) => {
              router.push(
                `#/workspaces/${encodeURIComponent(
                  workspaceId,
                )}/components/${encodeURIComponent(component.id)}/edit`,
              );
            },
          },
          {
            tooltip: 'Delete component',
            glyph: 'Delete',
            execute: async (component) => {
              await workspace.deleteComponent(component.id);
              window.location.reload();
            },
          },
        ]}
      />
    {/if}
  </Panel>
</Layout>
