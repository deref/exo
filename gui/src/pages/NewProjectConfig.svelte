<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import Spinner from '../components/Spinner.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';

  export let params = { starter: '' };

  const { starter } = params;

  let name = starter;
  let error: Error | null = null;

  let workingDirectory: string | null = null;

  const setWorkingDirectory = (dir: string) => {
    workingDirectory = dir;
  };

  (async () => {
    setWorkingDirectory(await api.kernel.getUserHomeDir());
  })();
</script>

<Layout>
  <CenterFormPanel title={`New project: ${starter}`} backRoute="#/new-project">
    <form
      on:submit|preventDefault={async () => {
        error = null;
        let workspaceId;
        // XXX Replace this old createWorkspace with templated createProject
        try {
          workspaceId = await api.kernel.createWorkspace(
            workingDirectory, // Currently this doesn't actually create the new directory, see XXX note above
          );
          router.push(`/workspaces/${encodeURIComponent(workspaceId)}`);
        } catch (ex) {
          if (!isClientError(ex)) {
            throw ex;
          }
          error = ex;
        }
        // XXX Hack to address lack of GUI for applying procfiles, etc.
        try {
          await api.workspace(workspaceId).apply();
        } catch (ex) {
          // Swallow error.
          console.error(ex);
        }
      }}
    >
      <h1>New project</h1>

      <label for="name">Name:</label>
      <Textbox
        bind:value={name}
        name="name"
        id="name"
        placeholder={starter}
        --input-width="100%"
      />
      <div style="height:32px" />
      {#if workingDirectory}
        <label for="root">Root:</label>
        <h2><span>{workingDirectory}<span>/{name}</span></span></h2>
        {#await api.kernel.readDir(workingDirectory)}
          <Spinner />
        {:then readdir}
          {#if readdir.parent !== null}
            <button
              on:click={() => setWorkingDirectory(String(readdir.parent?.path))}
            >
              ..
            </button>
          {/if}
          <div class="directories">
            {#each readdir.entries
              .filter((x) => x.isDirectory)
              .sort((x, y) => (x.name[0] > y.name[0] ? 1 : -1)) as entry}
              <button on:click={() => setWorkingDirectory(entry.path)}
                >{entry.name}</button
              >
            {/each}
          </div>
        {/await}
      {:else}
        <Spinner />
      {/if}
      <SubmitButton>Create project</SubmitButton>
    </form>
    <ErrorLabel value={error} />
  </CenterFormPanel>
</Layout>

<style>
  h2 span span {
    color: var(--grey-9-color);
  }

  .directories {
    margin: 12px 0;
  }

  .directories button {
    border-radius: 0;
  }

  .directories button:first-of-type {
    border-top-left-radius: 5px;
    border-top-right-radius: 5px;
  }

  .directories button:last-of-type {
    border-bottom-left-radius: 5px;
    border-bottom-right-radius: 5px;
  }

  button {
    background: var(--primary-bg-color);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 5px;
    padding: 5px 10px;
    position: relative;
    display: grid;
    width: 100%;
    grid-template-columns: max-content 2fr;
    align-items: center;
    gap: 12px;
    margin-top: 1px;
  }

  button:hover {
    background: var(--grey-e-color);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--grey-c-color);
    box-shadow: var(--button-active-shadow);
  }
</style>
