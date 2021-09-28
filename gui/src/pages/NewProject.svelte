<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api } from '../lib/api';
  import * as router from 'svelte-spa-router';

  let search: string = '';
</script>

<Layout>
  <CenterFormPanel title="New project" backRoute="/">
    <h1>New project</h1>

    <button
      on:click={() => {
        router.push('#/new-project/empty');
      }}
    >
      Empty project
    </button>

    <div class="search">
      <p>Select a starter for your new project:</p>
      <Textbox
        placeholder="Search..."
        bind:value={search}
        --input-width="100%"
        autofocus
      />
    </div>

    {#await api.kernel.describeTemplates()}
      <Spinner />
    {:then templates}
      {#each templates.filter((x) => x.displayName
            .toLocaleLowerCase()
            .slice(0, search.length) === search.toLocaleLowerCase()) as template}
        <button
          on:click={() => {
            router.push(`#/new-project/${template.name}`);
          }}
        >
          <Icon glyph={template.glyph} />
          {template.displayName}
        </button>
      {/each}
    {:catch error}
      <ErrorLabel value={error} />
    {/await}
  </CenterFormPanel>
</Layout>

<style>
  .search {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    align-items: center;
    margin-top: 24px;
  }

  .search :global(input) {
    height: 36px !important;
  }

  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 12px 18px;
    position: relative;
    display: grid;
    width: 100%;
    grid-template-columns: max-content 2fr;
    align-items: center;
    text-align: left;
    gap: 8px;
    margin-top: 8px;
  }

  button:not(:first-of-type) {
    padding: 8px 12px;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }
</style>
