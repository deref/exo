<script context="module" lang="ts">
  export interface Params {
    workspace: string;
  }
</script>

<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Button from '../components/Button.svelte';
  import Textbox from '../components/Textbox.svelte';
  import EditAs from '../components/form/EditAs.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import TextEditor from '../components/TextEditor.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import DockerSVG from '../components/mono/DockerSVG.svelte';
  import { api, isClientError } from '../lib/api';
  import { setLogVisibility } from '../lib/logs/visible-logs';
  import * as router from 'svelte-spa-router';

  export let params: Params;

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;
  const workspaceNewComponentRoute = `/workspaces/${encodeURIComponent(
    workspaceId,
  )}/new-component`;

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

  export let componentType: string;
  export let displayType: string;
</script>

<Layout>
  <CenterFormPanel
    title={`New ${displayType}`}
    backRoute={workspaceNewComponentRoute}
  >
    <h1><DockerSVG /> New {displayType}</h1>
    <form
      on:submit|preventDefault={async () => {
        try {
          const { id } = await workspace.createComponent(
            name,
            componentType,
            spec,
          );

          setLogVisibility(id, true);

          router.push(workspaceRoute);
        } catch (ex) {
          if (!isClientError(ex)) {
            throw ex;
          }
          error = ex;
        }
      }}
    >
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
      <div class="buttons">
        <Button type="submit">Create {displayType}</Button>
      </div>
      <div style="margin: 24px 0;">
        <ErrorLabel value={error} />
      </div>
    </form>
  </CenterFormPanel>
</Layout>

<style>
  details {
    margin: 24px 0;
  }

  summary {
    margin-bottom: 12px;
    user-select: none;
    cursor: pointer;
  }

  label {
    display: block;
    margin-bottom: 8px;
  }

  h1 {
    display: flex;
    align-items: center;
    gap: 18px;
    font-size: 24px;
    font-weight: 500;
    margin: 0;
    margin-bottom: 36px;
  }

  .buttons {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin-top: 8px;
  }
</style>
