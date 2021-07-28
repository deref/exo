<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { api } from '../lib/api';
  import Button from './Button.svelte';

  let installedVersion: string | null = null;
  let latestVersion: string | null = null;

  const doUpgrade = async () => {
    await api.kernel.upgrade();
  }

  let fetchTimeout = null;
  const refreshVersion = async () => {
    const { installed, latest, current } = await api.kernel.getVersion();
    installedVersion = installed;
    if (!current || true) {
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
    | Update available: <strong>{latestVersion}</strong> <Button on:click={doUpgrade}>Get It!</Button>
  {/if}
</section>
