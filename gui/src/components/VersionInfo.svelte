<script lang="ts">
  import { onDestroy } from 'svelte';
  import { api } from '../lib/api';
  import Button from './Button.svelte';
  import SpinnerSVG from './mono/spinner.svelte';

  let installedVersion: string | null = null;
  let latestVersion: string | null = null;
  let upgrading = false;

  const doUpgrade = async () => {
    upgrading = true;
    await api.kernel.upgrade();
  };

  let fetchTimeout: NodeJS.Timeout | null = null;
  const refreshVersion = async () => {
    const { installed, latest, current } = await api.kernel.getVersion();
    // The server just changed installed version - reload.
    if (installedVersion !== null && installedVersion != installed) {
      console.log({ installedVersion, installed });
      // window.location.reload();
    }

    installedVersion = installed;
    if (!current && latest !== undefined) {
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
    &nbsp;
    <Button size="small" on:click={doUpgrade}>
      {#if upgrading}
        Updating &nbsp;<SpinnerSVG />
      {:else}
        Upgrade
      {/if}
    </Button>
  {/if}
  {#if import.meta.env.MODE === 'development'}
    &nbsp;
    <strong>DEV MODE</strong>
  {/if}
</section>

<style>
  section {
    width: 48px;
    display: flex;
    align-items: center;
    padding: 12px 0;
    writing-mode: vertical-rl;
    text-orientation: mixed;
  }
</style>
