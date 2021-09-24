<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import Offline from './pages/Offline.svelte';
  import { isRunning, isAuthenticated } from './lib/global/server-status';
</script>

{#if $isRunning && $isAuthenticated}
  <Router {routes} />
{:else}
  <Offline />
{/if}
