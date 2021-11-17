<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import type { IconGlyph } from '../components/Icon.svelte';

  const queryString = router.querystring;
  const root = new URLSearchParams($queryString ?? '').get('root');

  let search: string = '';

  const logoGlyph = (ig: string) => (ig || 'Doc') as IconGlyph;
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

      {#await api.kernel.describeTemplates()}
        <Spinner />
      {:then templates}
        {#each templates
          .sort((_, y) => (y.iconGlyph ? 1 : -1))
          .filter((x) => x.displayName
                .toLocaleLowerCase()
                .slice(0, search.length) === search.toLocaleLowerCase()) as template}
          <button
            on:click={() => {
              router.push(
                `#/new-project/${template.name}${
                  !!root ? `?root=${encodeURIComponent(root)}` : ''
                }`,
              );
            }}
          >
            <Icon glyph={logoGlyph(template.iconGlyph)} />
            {template.displayName}
          </button>
        {/each}
      {:catch error}
        <ErrorLabel value={error} />
      {/await}
    </div>
  </CenterFormPanel>
</Layout>

<style>
  .search {
    display: grid;
    width: 100%;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    align-items: center;
    margin-top: 24px;
  }

  .search :global(input) {
    height: 36px !important;
  }

  button :global(svg) {
    height: 18px;
    width: 18px;
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
    gap: 12px;
    text-overflow: ellipsis;
    white-space: nowrap;
    overflow: hidden;
  }

  .search button {
    padding: 16px;
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
