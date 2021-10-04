<script context="module" lang="ts">
  export interface Params {
    workspace: string;
  }
</script>

<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import EditAs from '../components/form/EditAs.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import TextEditor from '../components/TextEditor.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api, isClientError } from '../lib/api';
  import { setLogVisibility } from '../lib/logs/visible-logs';
  import * as router from 'svelte-spa-router';

  export let params: Params;

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceComponentsRoute = `/workspaces/${encodeURIComponent(
    workspaceId,
  )}/components`;

  let error: Error | null = null;

  let mode = 'compose';
  const editorModes = [
    {
      id: 'compose',
      name: 'Compose YAML',
    },
    // {
    //   id: 'form',
    //   name: 'Form',
    // },
  ];

  let name: string = '';
  let spec: string = '';

  let componentType = 'process';
  let componentName = 'test-name';

  export let displayType: string;

  $: pageTitle = `Edit ${componentType} ${componentName}`;
</script>

<Layout>
  <CenterFormPanel title={pageTitle} backRoute={workspaceComponentsRoute}>
    <h1><Icon glyph="LogoDocker" />{pageTitle}</h1>
    <form on:submit|preventDefault={async () => {}}>
      <div class="group">
        <label for="name">Name:</label>
        <Textbox id="name" name="name" bind:value={name} --input-width="100%" />
      </div>

      <EditAs bind:mode {editorModes} />
      {#if mode === 'compose'}
        <div class="group">
          <label for="spec">Spec:</label>
          <TextEditor id="spec" bind:value={spec} language="yaml" />
        </div>
        <details>
          <summary>Show/hide example</summary>
          <slot />
        </details>
      {:else}
        <!-- GUI form edit mode -->
      {/if}
      <SubmitButton>Create {displayType}</SubmitButton>
      <div style="margin: 24px 0;">
        <ErrorLabel value={error} />
      </div>
    </form>
  </CenterFormPanel>
</Layout>
