<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { api } from '../lib/api';
  import Button from './Button.svelte';

  let installedVersion: string | null = null;
  let latestVersion: string | null = null;

  const doUpgrade = async () => {
    await api.kernel.upgrade();
  };

  let fetchTimeout = null;
  const refreshVersion = async () => {
    const { installed, latest, current } = await api.kernel.getVersion();
    // The server just changed version - reload.
    if (installedVersion !== null && installedVersion != latestVersion) {
      window.location.reload();
    }

    installedVersion = installed;
    if (!current) {
      latestVersion = latest;
    }
    fetchTimeout = setTimeout(refreshVersion, 60000);
  };
  refreshVersion();

  onDestroy(() => {
    if (fetchTimeout !== null) {
      clearTimeout(fetchTimeout);
    }
  });
</script>

<section>
  exo {installedVersion || ''}
  {#if latestVersion !== null}
    | Update available: <strong>{latestVersion}</strong>
    <Button on:click={doUpgrade}>Get It!</Button>
  {/if}
</section>

<style>
  section {
    font-size: 12px;
    color: #666;
  }
</style>
