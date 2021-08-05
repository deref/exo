<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import StringLabel from '../components/StringLabel.svelte';
  import ComponentTable from '../components/ComponentTable.svelte';
  import WithLeftWorkspaceNav from '../components/WithLeftWorkspaceNav.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
</script>

<Layout>
  <WithLeftWorkspaceNav {workspaceId} active="Components">
    <section>Components</section>
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
  </WithLeftWorkspaceNav>
</Layout>
