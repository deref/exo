<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
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
  <Panel --panel-padding="0">
    <div>
      <ProcessList {workspace} {workspaceId} />
      <LogsViewer {workspace} {workspaceId} />
    </div>
  </Panel>
</Layout>

<style>
  div {
    height: 100%;
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 30px;
    margin: 0 30px;
    padding-bottom: 16px;
  }
</style>
