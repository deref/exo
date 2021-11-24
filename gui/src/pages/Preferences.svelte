<script lang="ts">
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import IconButton from '../components/IconButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { theme, themeOptions } from '../lib/theme';
  import { api } from '../lib/api';
  import { modal } from '../lib/modal';
  import ModalTest from '../components/modal/ModalTest.svelte';

  const kernel = api.kernel;
  const makeRequest = () => kernel.getEsvUser('https://secrets.deref.io');
  let derefUser = makeRequest();

  const showModal1 = () => {
    modal.set(ModalTest);
  };
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
          {#await derefUser}
            Loading...
          {:then user}
            {#if user}
              <p>User: {user.email}</p>
              <div class="button-row">
                <Button
                  on:click={async () => {
                    await kernel.unauthEsv();
                    derefUser = makeRequest();
                  }}
                >
                  Unauthenticate</Button
                >
              </div>
            {:else}
              Not logged in
            {/if}
          {/await}
        </div>
      </div>
      <section>
        <h2>Test modal</h2>

        <Button on:click={showModal1}>Show modal 1</Button>
      </section>
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
