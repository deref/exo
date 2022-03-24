<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import OfflinePage from './components/OfflinePage.svelte';
  import { initClient as initGraphqlClient, online } from './lib/graphql';

  initGraphqlClient();
</script>

{#if $online}
  <Router {routes} />
{:else}
  <OfflinePage />
{/if}
