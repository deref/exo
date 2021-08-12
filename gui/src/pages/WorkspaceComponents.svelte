<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import MonoPanel from '../components/MonoPanel.svelte';
  import StringLabel from '../components/StringLabel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import ComponentTable from '../components/ComponentTable.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Components" slot="navbar" />
  <MonoPanel>
    <h1>Components</h1>
    <ComponentTable
      load={workspace.describeComponents}
      columns={[
        {
          title: 'id',
          component: StringLabel,
          getValue: (component) => component.id,
        },
        {
          title: 'name',
          component: StringLabel,
          getValue: (component) => component.name,
        },
        {
          title: 'type',
          component: StringLabel,
          getValue: (component) => component.type,
        },
      ]}
    />
  </MonoPanel>
</Layout>

<style>
  h1 {
    padding: 0;
    margin: 0;
    margin-bottom: 30px;
  }
</style>
