<script>
  import { api } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';

  import Code from '../components/Code.svelte';

  let pingInterval = null;
  onMount(() => {
    pingInterval = setInterval(() => api.kernel.ping(), 1000);
  });
  onDestroy(() => {
    clearInterval(pingInterval);
  });
</script>

<main>
  <div>
    <h3>
      <code>exo</code> server is offline
    </h3>
    <p class="hint">
      Run <Code>exo daemon</Code> in your terminal to start.
    </p>
  </div>
</main>

<style>
  main {
    display: flex;
    align-items: center;
    justify-content: center;
    font-family: Inter, sans-serif;
    font-size: 18px;
    position: absolute;
    height: 100%;
    width: 100%;
  }
  p,
  h3 {
    text-align: center;
  }
</style>
