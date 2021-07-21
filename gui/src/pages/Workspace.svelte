<script lang="ts">
  import ProcessList from '../components/ProcessList.svelte'
  import LogsViewer from '../components/LogsViewer.svelte'
  import { api } from '../lib/api';
  
  export let params = { workspace: '' };
  
  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
</script>

<div class="container">
  <ProcessList workspace={workspace} workspaceId={workspaceId}/>
  <LogsViewer workspace={workspace}/>
</div>

<style>

  .container {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 30px;
    margin: 0 30px;
  }

</style>
