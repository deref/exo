<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import TwoColumn from '../components/TwoColumn.svelte';
  import LogsViewer from '../components/LogsViewer.svelte';
  import ProcessList from '../components/ProcessList.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <TwoColumn>
    <ProcessList slot="left" {workspaceId} />
    <LogsViewer slot="right" {workspaceId} />
  </TwoColumn>
</Layout>
