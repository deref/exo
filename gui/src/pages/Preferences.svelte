<script lang="ts">
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import IconButton from '../components/IconButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { theme, themeOptions } from '../lib/theme';
  import { api } from '../lib/api';

  const kernel = api.kernel;
  // FIXME: this should point at prod.
  const derefUser = kernel.getEsvUser('http://localhost:5000');
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
      <div>
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
        <div class="group-header">
          <h2>Deref</h2>
        </div>
        <div>
          {#await derefUser then user}
            {#if user}
              Email: {user.email}
            {:else}
              Not logged in
            {/if}
          {/await}
        </div>
      </div>
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
    margin-bottom: 2em;
  }
</style>
