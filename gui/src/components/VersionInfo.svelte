<script lang="ts">
  import Icon from './Icon.svelte';
  import Button from './Button.svelte';
  import Spinner from './Spinner.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import { onDestroy } from 'svelte';
  import { api } from '../lib/api';
  import { nonNull } from '../lib/util';
  import { subscribe } from 'svelte-apollo';
  import gql from 'graphql-tag';

  import { modal } from '../lib/modal';
  import { bind } from './modal/Modal.svelte';
  import ModalDefaultPopup from './modal/ModalDefaultPopup.svelte';

  type TODO_QUERY_DATA = {
    system: {
      version: {
        installed: string;
        managed: boolean;
        upgrade: string | null;
      };
    };
  }; // XXX
  const q = subscribe<TODO_QUERY_DATA>(gql`
    subscription {
      system: systemChange {
        version {
          installed
          managed
          upgrade
        }
      }
    }
  `);
  $: version = $q.data?.system?.version;

  let upgrading = false;
  const doUpgrade = async () => {
    if (nonNull(version).managed) {
      modal.set(
        bind(ModalDefaultPopup, {
          title: 'Cannot auto-update.',
          message: `This upgrade procedure only supports installation of exo that were performed with the exo install script.\n\nPlease upgrade exo with your package manager or by uninstalling and reinstalling with the official exo upgrade script.`,
        }),
      );
      return;
    }
    upgrading = true;
    q.subscribe(() => {
      // If an event comes in now, assume it means a new version was installed.
      window.location.reload();
    });
    await api.kernel.upgrade();
  };
</script>

<div class="dropdown-wrapper">
  <NavbarButton>
    {#if version?.upgrade}
      <div class="upgrade-available" />
    {:else}
      <Icon glyph="Ellipsis" />
    {/if}
  </NavbarButton>
  <div class="dropdown version">
    <section>
      <div>
        exo {version?.installed ?? ''}
        {#if import.meta.env.MODE === 'development'}
          <strong>dev</strong>
        {/if}
      </div>
      {#if version?.upgrade}
        <div>Update available: <strong>{version.upgrade}</strong></div>
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
