<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import Spinner from '../components/Spinner.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import DirectoryBrowser from '../components/DirectoryBrowser.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import type { ReadDirResult } from '../lib/api';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';

  export let params = { starter: '' };

  const { starter } = params;

  let name = starter;
  let templateUrl: string | null = null;
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
</script>

<Layout>
  <CenterFormPanel title={`New ${starter} project`} backRoute="#/new-project">
    <form
      on:submit|preventDefault={async () => {
        error = null;
        let workspaceId;

        if (starter !== 'empty') {
          templateUrl = String(
            (await api.kernel.describeTemplates()).find(
              (x) => x.name === starter,
            )?.url,
          );
        }

        try {
          workspaceId = await api.kernel.createProject(
            `${withSlash(String(workingDirectory))}${name}`,
            templateUrl,
          );
          router.push(`/workspaces/${encodeURIComponent(workspaceId)}`);
        } catch (ex) {
          if (!(ex instanceof Error) || !isClientError(ex)) {
            throw ex;
          }
          error = ex;
        }
      }}
    >
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
          <DirectoryBrowser {dir} handleClick={setWorkingDirectory} />
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
</style>
