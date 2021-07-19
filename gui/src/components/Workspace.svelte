<script lang="ts">
  import ProcessList from './ProcessList.svelte'
  import LogsViewer from './LogsViewer.svelte'
  import { api } from '../lib/api';
  import { workspaceId } from '../lib/workspaces/store';
  
  const workspace = api.workspace($workspaceId);
</script>

<div class="layout">
  <ProcessList workspace={workspace}/>
  <LogsViewer workspace={workspace}/>
</div>

<style>

  .layout {
    display: grid;
    grid-template-columns: 2fr 3fr;
    gap: 30px;
    margin: 0 30px;
  }

</style>
