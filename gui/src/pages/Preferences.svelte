<script lang="ts">
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import IconButton from '../components/IconButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { theme, themeOptions } from '../lib/theme';

  let typographyPrefs = [
    {
      variable: 'main-font-size',
      value: '16px',
    },
    {
      variable: 'log-font-size',
      value: '15px',
    },
    {
      variable: 'ligatures-logs',
      value: 'none',
    },
    {
      variable: 'ligatures-code',
      value: 'none',
    },
  ];
</script>

<Layout>
  <CenterFormPanel title="Preferences">
    <div slot="actions">
      <IconButton
        glyph="Reset"
        tooltip="Reset to defaults"
        on:click={() => theme.apply('auto')}
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

      <div class="group">
        <div class="group-header">
          <h2>Typography</h2>
        </div>
        {#each typographyPrefs as pref}
          <div class="input-row">
            <code>{pref.variable}</code>
            <Textbox bind:value={pref.value} --input-width="100%" />
          </div>
        {/each}
      </div>
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
