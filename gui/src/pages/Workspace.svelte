<script lang="ts">
  import Layout from '../components/Layout.svelte';
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
  <div>
    <ProcessList {workspace} {workspaceId} />
    <LogsViewer {workspace} {workspaceId} />
  </div>
</Layout>

<style>
  div {
    height: 100%;
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 1px;
  }
</style>
