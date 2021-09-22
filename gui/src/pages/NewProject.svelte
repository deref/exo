<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import * as router from 'svelte-spa-router';

  const api = {
    kernel: {
      describeTemplates: async () => [
        {
          name: 'empty-project',
          displayName: 'Empty project',
          url: 'https://github.com/deref/exo-starters/empty-project',
        },
        {
          name: 'nextjs-prisma',
          displayName: 'NextJS Prisma',
          url: 'https://github.com/railwayapp/starters/tree/master/examples/nextjs-prisma',
        },
        {
          name: 'http-nodejs',
          displayName: 'HTTP Module',
          url: 'https://github.com/railwayapp/starters/tree/master/examples/http-nodejs',
        },
        {
          name: 'laravel',
          displayName: 'Laravel Starter',
          url: 'https://github.com/railwayapp/starters/tree/master/examples/laravel',
        },
      ],
    },
  };
</script>

<Layout>
  <CenterFormPanel title="New project" backRoute="/">
    <h1>New project</h1>
    <p>Select a starter for your new project:</p>
    {#await api.kernel.describeTemplates()}
      <Spinner />
    {:then templates}
      {#each templates as template}
        <button
          on:click={() => {
            router.push(`#/new-project/${template.name}`);
          }}>{template.displayName}</button
        >
      {/each}
    {/await}
  </CenterFormPanel>
</Layout>

<style>
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
    gap: 8px;
    margin-top: 8px;
  }

  button > b {
    max-width: 6em;
  }

  button > * {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
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
