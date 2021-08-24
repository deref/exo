<script lang="ts">
  import Button from './Button.svelte';
  import Spinner from './Spinner.svelte';
  import EllipsisSVG from './mono/EllipsisSVG.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import { onDestroy } from 'svelte';
  import { api } from '../lib/api';

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

<div class="dropdown-wrapper">
  <NavbarButton>
    {#if latestVersion !== null}
      <div class="upgrade-available" />
    {:else}
      <EllipsisSVG />
    {/if}
  </NavbarButton>
  <div class="dropdown version">
    <section>
      <div>exo {installedVersion || ''}</div>
      {#if latestVersion !== null}
        <div>Update available: <strong>{latestVersion}</strong></div>
        <Button small on:click={doUpgrade}>
          {#if upgrading}
            Updating &nbsp;<Spinner inline />
          {:else}
            Upgrade
          {/if}
        </Button>
      {/if}
    </section>
  </div>
</div>

<style>
  .upgrade-available {
    width: 20px;
    height: 20px;
    border-radius: 10px;
    background: var(--error-color);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .upgrade-available::after {
    content: '1';
    color: white;
    font-size: 12px;
    font-weight: 500;
  }

  section:nth-child(2) {
    margin-top: 8px;
  }

  .dropdown-wrapper:not(:hover) .dropdown {
    display: none;
  }

  .dropdown {
    position: absolute;
    bottom: 8px;
    left: calc(100% - 8px);
    border-radius: 5px;
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
  }

  .dropdown.version {
    font-size: 15px;
    padding: 12px 18px;
    min-width: 240px;
  }
</style>
