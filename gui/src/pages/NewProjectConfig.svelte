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

  const queryString = router.querystring;
  const root = new URLSearchParams($queryString ?? '').get('root');

  let name = starter;
  let error: Error | null = null;

  let workingDirectory: string | null = null;
  let homeDirectory: string = '/';

  const setWorkingDirectory = (dir: string) => {
    workingDirectory = dir;
    dirPromise = api.kernel.readDir(workingDirectory);
  };

  let dirPromise: Promise<ReadDirResult> = api.kernel
    .getUserHomeDir()
    .then((dir) => {
      if (root) {
        dir = root;
      }
      setWorkingDirectory(dir);
      homeDirectory = dir;
      return dirPromise;
    });

  const withSlash = (wd: string) => (wd.endsWith('/') ? wd : wd + '/');

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
      if (!isClientError(ex)) {
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
          <DirectoryBrowser
            {dir}
            homePath={homeDirectory}
            handleClick={setWorkingDirectory}
          />
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
