<script lang="ts">
  import Icon from './Icon.svelte';
  import Button from './Button.svelte';
  import Spinner from './Spinner.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import { onDestroy } from 'svelte';
  import { api } from '../lib/api';

  import { modal } from '../lib/modal';
  import { bind } from './modal/Modal.svelte';
  import ModalDefaultPopup from './modal/ModalDefaultPopup.svelte';

  let installedVersion: string | null = null;
  let latestVersion: string | null = null;
  let isManaged: boolean = false;
  let upgrading = false;

  const showUpdateInstallMethodModal = () => {
    modal.set(
      bind(ModalDefaultPopup, {
        title: 'Cannot auto-update.',
        message: `This upgrade procedure only supports installation of exo that were performed with the exo install script.\n\nPlease upgrade exo with your package manager or by uninstalling and reinstalling with the official exo upgrade script.`,
      }),
    );
  };

  const doUpgrade = async () => {
    if (isManaged) {
      showUpdateInstallMethodModal();
      return;
    }
    upgrading = true;
    await api.kernel.upgrade();
  };

  let fetchTimeout: NodeJS.Timeout | null = null;
  const refreshVersion = async () => {
    const { installed, latest, current, managed } =
      await api.kernel.getVersion();
    // The server just changed installed version - reload.
    if (installedVersion !== null && installedVersion != installed) {
      console.log({ installedVersion, installed });
      // window.location.reload();
    }

    isManaged = managed;

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
      <Icon glyph="Ellipsis" />
    {/if}
  </NavbarButton>
  <div class="dropdown version">
    <section>
      <div>
        exo {installedVersion || ''}
        {#if import.meta.env.MODE === 'development'}
          <strong>dev</strong>
        {/if}
      </div>
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

  .dropdown-wrapper:not(:hover):not(:focus-within) .dropdown {
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
