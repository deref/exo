<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import MonoPanel from '../components/MonoPanel.svelte';
</script>

<Layout>
  <MonoPanel title="Error" backRoute="/">
    <div>
      <h1>404 page not found</h1>
    </div>
  </MonoPanel>
</Layout>

<style>
  div {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
</style>
