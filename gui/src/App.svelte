<script lang="ts">
  import ProcessList from './components/ProcessList.svelte'
  import LogsViewer from './components/LogsViewer.svelte'
</script>

<main>
  <header>
    <img src="/deref-rounded-icon.png" alt="Deref" height="24px" /> <h1>Exo</h1>
  </header>


  <div class="layout">
    <ProcessList />
    <LogsViewer />
  </div>
</main>

<style>

  header {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px;
    box-shadow: 0px 6px 9px -6px #00000022, 0px 0.25px 0px 1px #00000022;
  }
  
  main {
    /* */
  }

  h1 {
    font-size: 20px;
    font-weight: 550;
    margin: 0;
  }

  .layout {
    display: grid;
    grid-template-columns: 2fr 3fr;
    gap: 30px;
    margin: 0 30px;
  }

</style>
