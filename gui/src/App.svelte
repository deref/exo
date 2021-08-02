<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import Offline from './pages/Offline.svelte';
  import { isRunning } from './lib/global/server-status';
</script>


{#if $isRunning}
  <Router {routes} />
{:else}
  <Offline />
{/if}
