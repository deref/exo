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
  <WithLeftWorkspaceNav {workspaceId} active="Networking">
    <section>Networks</section>
    <ComponentTable
      load={workspace.describeNetworks}
      columns={[
        {
          title: 'id',
          component: StringLabel,
          getValue: (network) => network.id,
        },
        {
          title: 'name',
          component: StringLabel,
          getValue: (network) => network.name,
        },
      ]}
    />
  </WithLeftWorkspaceNav>
</Layout>
