<script lang="ts">
  import Code from '../components/Code.svelte';
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import { api } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';

  let pingInterval: null | ReturnType<typeof setInterval> = null;
  onMount(() => {
    pingInterval = setInterval(() => api.kernel.ping(), 1000);
  });
  onDestroy(() => {
    if (pingInterval) {
      clearInterval(pingInterval);
    }
  });
</script>

<Layout>
  <Panel>
    <div>
      <h3>
        <Code>exo</Code> server is offline
      </h3>
      <p>
        Run <Code>exo daemon</Code> in your terminal to start.
      </p>
    </div>
  </Panel>
</Layout>

<style>
  div {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: 18px;
    height: 100%;
  }
</style>
