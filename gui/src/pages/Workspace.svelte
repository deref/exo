<script lang="ts">
  import ProcessList from '../components/ProcessList.svelte'
  import LogsViewer from '../components/LogsViewer.svelte'
  import { api } from '../lib/api';
  
  export let params = { workspace: '' };
  
  const workspace = api.workspace(params.workspace);
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
