<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import Offline from './pages/Offline.svelte';
  import { isRunning, isAuthenticated } from './lib/global/server-status';
  import { ApolloClient, InMemoryCache } from '@apollo/client';
  import { setClient } from 'svelte-apollo';

  const apiPort = parseInt(import.meta.env.VITE_API_PORT as string);
  const client = new ApolloClient({
    cache: new InMemoryCache(),
    uri: `http://localhost:${apiPort}/_exo/graphql`,
  });
  setClient(client);
</script>

{#if $isRunning && $isAuthenticated}
  <Router {routes} />
{:else}
  <Offline />
{/if}
