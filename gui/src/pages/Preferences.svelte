<script lang="ts">
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import IconButton from '../components/IconButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { theme, themeOptions } from '../lib/theme';
  import { preferences } from '../lib/preferences';
  import { onMount } from 'svelte';

  let prefGroups = [
    {
      title: 'Typography',
      prefs: [
        {
          name: 'main-font-size',
        },
        {
          name: 'log-font-size',
        },
        {
          name: 'ligatures-logs',
        },
        {
          name: 'ligatures-code',
        },
      ],
    },
  ];

  type Preferences = Record<string, string>;

  let dirtyPrefs: Preferences = {};

  onMount(() => {
    dirtyPrefs = $preferences;
  });
</script>

<Layout>
  <CenterFormPanel title="Preferences">
    <div slot="actions">
      <IconButton
        glyph="Reset"
        tooltip="Reset to defaults"
        on:click={() => {
          theme.apply('auto');
          preferences.reset();
          setTimeout(() => {
            dirtyPrefs = $preferences;
          }, 50);
        }}
      />
    </div>
    <div>
      <div class="group">
        <div class="group-header">
          <h2>Theme &amp; GUI</h2>
        </div>
        <div class="button-row">
          {#each themeOptions as option}
            <Button
              on:click={() => theme.apply(option)}
              inset={$theme === option}
            >
              <code>{option}</code>
            </Button>
          {/each}
        </div>
      </div>

      {#each prefGroups as group}
        <div class="group">
          <div class="group-header">
            <h2>{group.title}</h2>
          </div>
          {#each group.prefs as pref}
            <div class="input-row">
              <code>{pref.name}</code>
              <Textbox
                bind:value={dirtyPrefs[pref.name]}
                on:input={() => preferences.apply({ ...dirtyPrefs })}
                --input-width="100%"
              />
            </div>
          {/each}
        </div>
      {/each}
    </div>
  </CenterFormPanel>
</Layout>

<style>
  .input-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    align-items: center;
    margin-bottom: 12px;
  }

  .group:not(:last-child) {
    margin-bottom: 40px;
  }

  .group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .button-row {
    display: grid;
    grid-auto-flow: column;
    gap: 12px;
  }
</style>
