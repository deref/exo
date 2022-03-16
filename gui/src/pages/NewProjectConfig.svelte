<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import DirectoryBrowser from '../components/DirectoryBrowser.svelte';
  import type { DirectoryStore } from '../components/DirectoryBrowser.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import { query } from '../lib/graphql';
  import { derived, writable } from 'svelte/store';
  import { nonNull } from '../lib/util';

  export let params = { starter: '' };

  const { starter } = params;
  const templateName = starter ?? null;

  const queryString = router.querystring;
  const workspace = new URLSearchParams($queryString ?? '').get('workspace');

  let name = templateName;
  let formError: Error | null = null;

  // Null implies home.
  let directoryPath = writable<string | null>(workspace);

  const q = query(
    `#graphql
    query ($directoryPath: String) {
      fileSystem {
        homePath
        file: fileOrHome(path: $directoryPath) {
          path
          parentPath
          children {
            name
            path
            isDirectory
          }
        }
      }
    }`,
    {
      variables: {
        directoryPath: $directoryPath,
      },
    },
  );

  const store: DirectoryStore = {
    ...derived(q, ($q) =>
      $q.loading || $q.error || !$q.data
        ? {
            ready: false as const,
            error: $q.error ?? null,
          }
        : {
            ready: true as const,
            homePath: nonNull($q.data).fileSystem.homePath,
            directory: nonNull(nonNull($q.data).fileSystem.file),
          },
    ),
    setDirectory: (value: string) => {
      directoryPath.set(value);
    },
  };

  $: homePath = $q.data?.fileSystem.homePath;

  const submitForm = async () => {
    formError = null;

    // Shouldn't really be possible since the form shouldn't be visible if the
    // working directory isn't set.
    if (!$store.ready) {
      formError = new Error('Directory not set');
      return;
    }

    try {
      const workspaceId = await api.kernel.createProject(
        `${$store.directory.path}${name}`,
        templateName, // XXX used to be templateUrl, handle server side.
      );
      await router.push(`/workspaces/${encodeURIComponent(workspaceId)}`);
    } catch (ex) {
      if (!isClientError(ex)) {
        throw ex;
      }
      formError = ex;
    }
  };
</script>

<Layout loading={$q.loading} error={null}>
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
      {#if homePath}
        <label for="root">Root:</label>
        <h2><span>{$directoryPath || homePath}<span>{name}</span></span></h2>
      {/if}
      <DirectoryBrowser {store} autofocus />
      <SubmitButton>Create project</SubmitButton>
    </form>
    <ErrorLabel value={formError} />
  </CenterFormPanel>
</Layout>

<style>
  h2 span span {
    color: var(--grey-9-color);
  }
</style>
