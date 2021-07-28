<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import ProcessList from '../components/ProcessList.svelte';
  import LogsViewer from '../components/LogsViewer.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
</script>

<Layout>
  <section>
    <ProcessList {workspace} {workspaceId} />
    <LogsViewer {workspace} {workspaceId} />
  </section>
</Layout>

<style>
  section {
    height: 100%;
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 30px;
    margin: 0 30px;
    padding-bottom: 16px;
  }
</style>
