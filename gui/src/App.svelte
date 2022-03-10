<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import Offline from './pages/Offline.svelte';
  import { isRunning, isAuthenticated } from './lib/global/server-status';
  import { initClient as initGraphqlClient } from './lib/graphql';

  initGraphqlClient();
</script>

{#if $isRunning && $isAuthenticated}
  <Router {routes} />
{:else}
  <Offline />
{/if}
