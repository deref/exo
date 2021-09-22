<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api } from '../lib/api';
  import * as router from 'svelte-spa-router';
</script>

<Layout>
  <CenterFormPanel title="New project" backRoute="/">
    <h1>New project</h1>
    <p>Select a starter for your new project:</p>

    <button
      on:click={() => {
        router.push('#/new-project/empty');
      }}
    >
      Empty project
    </button>

    {#await api.kernel.describeTemplates()}
      <Spinner />
    {:then templates}
      {#each templates as template}
        <button
          on:click={() => {
            router.push(`#/new-project/${template.name}`);
          }}
        >
          {template.displayName}
        </button>
      {/each}
    {:catch error}
      <ErrorLabel value={error} />
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

  button:not(:first-of-type) {
    padding: 6px 18px;
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
