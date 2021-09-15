<script lang="ts">
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import IconButton from '../components/IconButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import PreferencesSectionsConfig from '../components/form/PreferencesSectionsConfig.svelte';
  import ResetSVG from '../components/mono/ResetSVG.svelte';
  import { theme } from '../lib/theme';

  const resetAllPreferences = () => {
    theme.apply('auto');
  };
</script>

<Layout>
  <CenterFormPanel title="Preferences">
    <div slot="actions">
      <IconButton tooltip="Reset to defaults" on:click={resetAllPreferences}>
        <ResetSVG />
      </IconButton>
    </div>
    <div>
      <div class="group">
        <div class="group-header">
          <h2>Theme &amp; GUI</h2>
        </div>
        <div class="button-row">
          <Button
            on:click={() => {
              theme.apply('auto');
            }}
            inset={$theme === 'auto'}
          >
            Auto
          </Button>
          <Button
            on:click={() => {
              theme.apply('light');
            }}
            inset={$theme === 'light'}
          >
            Light
          </Button>
          <Button
            on:click={() => {
              theme.apply('dark');
            }}
            inset={$theme === 'dark'}
          >
            Dark
          </Button>
          <Button
            on:click={() => {
              theme.apply('black');
            }}
            inset={$theme === 'black'}
          >
            Black
          </Button>
        </div>
      </div>
      <!-- <PreferencesSectionsConfig /> -->
    </div>
  </CenterFormPanel>
</Layout>

<style>
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
