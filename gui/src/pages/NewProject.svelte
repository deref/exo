<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import * as router from 'svelte-spa-router';
  import { api } from '../lib/api';

  const starters = api.kernel.describeTemplates();
</script>

<Layout>
  <CenterFormPanel title="New project" backRoute="/">
    <h1>New project</h1>
    <p>Select a starter for your new project:</p>
    {#await starters}
      loading templates...
    {:then starters}
      {#each starters as starter}
        <button
          on:click={() => {
            router.push(`#/new-project/${starter.name}`);
          }}>{starter.displayName}</button
        >
      {/each}
    {:catch error}
      <p style="color: red">{error.message}</p>
    {/await}
  </CenterFormPanel>
</Layout>

<style>
  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 16px 24px;
    position: relative;
    display: grid;
    width: 100%;
    grid-template-columns: max-content 2fr;
    align-items: center;
    gap: 12px;
    margin-top: 12px;
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
