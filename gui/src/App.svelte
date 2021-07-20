<script lang="ts">
  import Workspace from './components/Workspace.svelte'
  import CreateWorkspacePage from './components/CreateWorkspacePage.svelte'
  import { workspaceId } from './lib/workspaces/store';
  
  // XXX Janky routing.
  console.log(window.location.search);
  const match = window.location.pathname.match(/\/workspaces\/(.*)/);
  if (match) {
    workspaceId.set(match[1]);
  }
</script>

<main>
  <header>
    <img src="/deref-rounded-icon.png" alt="Deref" height="24px" /> <h1>Exo</h1>
  </header>

  {#if $workspaceId}
    <Workspace/>
  {:else}
    <CreateWorkspacePage/>
  {/if}
</main>

<style>

  header {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px;
    box-shadow: 0px 6px 9px -6px #00000022, 0px 0.25px 0px 1px #00000022;
  }
  
  h1 {
    font-size: 20px;
    font-weight: 550;
    margin: 0;
  }

</style>
