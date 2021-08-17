<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import TwoColumn from '../components/TwoColumn.svelte';
  import LogsViewer from '../components/LogsViewer.svelte';
  import ProcessList from '../components/ProcessList.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <TwoColumn>
    <ProcessList slot="left" {workspace} {workspaceId} />
    <LogsViewer slot="right" {workspace} />
  </TwoColumn>
</Layout>
