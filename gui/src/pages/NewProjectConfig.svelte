<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import Spinner from '../components/Spinner.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import type { ReadDirResult } from '../lib/api';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';

  export let params = { starter: '' };

  const { starter } = params;

  let name = starter;
  let error: Error | null = null;

  let workingDirectory: string | null = null;

  const setWorkingDirectory = (dir: string) => {
    workingDirectory = dir;
    dirPromise = api.kernel.readDir(workingDirectory);
  };

  let dirPromise: Promise<ReadDirResult> = api.kernel
    .getUserHomeDir()
    .then((dir) => {
      setWorkingDirectory(dir);
      return dirPromise;
    });

  const withSlash = (wd: string) => (wd.slice(-1) === '/' ? wd : wd + '/');

  const submitForm = async () => {
    error = null;

    let templateUrl: string | null = null;
    if (starter !== 'empty') {
      const templates = await api.kernel.describeTemplates();
      const template = templates.find((x) => x.name === starter);
      if (!template) {
        error = new Error(`Could not find template with name ${starter}`);
        return;
      }
      templateUrl = template.url;
    }

    // Shouldn't really be possible since the form shouldn't be visible if the
    // working directory isn't set.
    if (!workingDirectory) {
      error = new Error('Working directory not set');
      return;
    }

    try {
      const workspaceId = await api.kernel.createProject(
        `${withSlash(workingDirectory)}${name}`,
        templateUrl,
      );
      await router.push(`/workspaces/${encodeURIComponent(workspaceId)}`);
    } catch (ex) {
      if (!(ex instanceof Error) || !isClientError(ex)) {
        throw ex;
      }
      error = ex;
    }
  };
</script>

<Layout>
  <CenterFormPanel title={`New ${starter} project`} backRoute="#/new-project">
    <form on:submit|preventDefault={submitForm}>
      <h1>New {starter} project</h1>

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
        <h2><span>{withSlash(workingDirectory)}<span>{name}</span></span></h2>
        {#await dirPromise}
          <Spinner />
        {:then dir}
          {#if dir.parent !== null}
            <button
              on:click={() => setWorkingDirectory(String(dir.parent?.path))}
            >
              ..
            </button>
          {:else}
            <button disabled> .. </button>
          {/if}
          <div class="directories">
            {#each dir.entries
              .filter((x) => x.isDirectory)
              .sort((x, y) => x.name.localeCompare(y.name)) as entry}
              <button on:click={() => setWorkingDirectory(entry.path)}
                >{entry.name}</button
              >
            {/each}
          </div>
        {:catch awaitError}
          <ErrorLabel value={awaitError} />
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

  button:disabled {
    cursor: default;
    opacity: 0.5;
  }

  button:not(:disabled):hover {
    background: var(--grey-e-color);
    box-shadow: var(--button-hover-shadow);
  }

  button:not(:disabled):active {
    background: var(--grey-c-color);
    box-shadow: var(--button-active-shadow);
  }
</style>
