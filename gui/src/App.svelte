<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import Offline from './pages/Offline.svelte';
  import { initClient as initGraphqlClient, online } from './lib/graphql';

  initGraphqlClient();
</script>

{#if $online}
  <Router {routes} />
{:else}
  <Offline />
{/if}
