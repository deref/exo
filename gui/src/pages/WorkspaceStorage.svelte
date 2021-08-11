<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import StringLabel from '../components/StringLabel.svelte';
  import ComponentTable from '../components/ComponentTable.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import MonoPanel from '../components/MonoPanel.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Storage" slot="navbar" />
  <MonoPanel>
    <section>
      <h1>Volumes</h1>
      <ComponentTable
        load={workspace.describeVolumes}
        columns={[
          {
            title: 'id',
            component: StringLabel,
            getValue: (volume) => volume.id,
          },
          {
            title: 'name',
            component: StringLabel,
            getValue: (volume) => volume.name,
          },
        ]}
      />
    </section>
  </MonoPanel>
</Layout>

<style>
  section {
    padding: 30px;
  }

  h1 {
    padding: 0;
    margin: 0;
    margin-bottom: 30px;
  }
</style>
