<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
</script>

<Layout>
  <Panel title="Error" backUrl="/">
    <div>
      <h1>404 page not found</h1>
    </div>
  </Panel>
</Layout>

<style>
  div {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
