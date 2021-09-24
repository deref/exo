<script lang="ts">
  import routes from './routes';
  import Router from 'svelte-spa-router';
  import Offline from './pages/Offline.svelte';
  import { isRunning } from './lib/global/server-status';

  const params = new URLSearchParams(window.location.search);
  const urlToken = params.get('token');
  if (urlToken) {
    document.cookie = `token=${urlToken}`;
    params.delete('token');
    window.location.search = params.toString();
  }
</script>

{#if $isRunning}
  <Router {routes} />
{:else}
  <Offline />
{/if}
