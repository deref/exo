<script lang="ts">
  import type { ReadDirResult } from '../lib/api';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api } from '../lib/api';

  let workingDirectory = '/';
  const setWorkingDirectory = (dir: string) => {
    workingDirectory = dir;
    currentPromise = api.kernel.readDir(workingDirectory);
  };
  let currentPromise: Promise<ReadDirResult> = api.kernel
    .getUserHomeDir()
    .then((dir) => {
      setWorkingDirectory(dir);
      return currentPromise;
    });

  export let params = { starter: '' };

  const { starter } = params;

  let name = starter;
  let error = 'ee';
</script>

<Layout>
  <CenterFormPanel title={`New project: ${starter}`} backRoute="#/new-project">
    <form on:submit|preventDefault={async () => {}}>
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
      <label for="root">Root:</label>
      <h2>
        <span
          >{workingDirectory.slice(-1) === '/'
            ? workingDirectory
            : workingDirectory + '/'}<span>{name}</span></span
        >
      </h2>
      <div>
        {#await currentPromise then current}
          <button
            disabled={workingDirectory === '/'}
            on:click={async () => {
              if (current.parent) {
                setWorkingDirectory(current.parent.path);
              }
            }}>..</button
          >
          <div class="directories">
            {#each current.entries as entry}
              {#if entry.isDirectory}
                <button
                  on:click={() => {
                    workingDirectory = entry.path;
                    currentPromise = api.kernel.readDir(workingDirectory);
                  }}>{entry.name}</button
                >
              {/if}
            {/each}
          </div>
        {:catch error}
          <!-- TODO better error handling -->
          <p>Error {error}</p>
        {/await}
      </div>
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

  button > b {
    max-width: 6em;
  }

  button > * {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
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
