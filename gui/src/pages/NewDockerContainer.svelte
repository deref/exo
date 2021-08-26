<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';
  import Button from '../components/Button.svelte';
  import Textbox from '../components/Textbox.svelte';
  import EditAs from '../components/form/EditAs.svelte';
  import CodeBlock from '../components/CodeBlock.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import TextEditor from '../components/TextEditor.svelte';
  import DockerSVG from '../components/mono/DockerSVG.svelte';
  import { api, isClientError } from '../lib/api';
  import { setLogVisibility } from '../lib/logs/visible-logs';
  import * as router from 'svelte-spa-router';

  export let params = { workspace: '' };

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

  let specExample = `image: postgres:9.4
environment:
    POSTGRES_USER: "postgres"
    POSTGRES_PASSWORD: "postgres"`;
</script>

<Layout>
  <Panel
    title="New Container"
    backRoute={workspaceNewComponentRoute}
    --panel-padding="2rem"
    --panel-overflow-y="scroll"
  >
    <div class="center-form">
      <h1><DockerSVG /> New Container</h1>
      <form
        on:submit|preventDefault={async () => {
          try {
            const { id } = await workspace.createComponent(
              name,
              'container',
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
          <Textbox
            id="name"
            name="name"
            bind:value={name}
            --input-width="100%"
          />
        </div>

        <EditAs bind:mode {editorModes} />
        {#if mode === 'compose'}
          <div class="group">
            <label for="spec">Spec:</label>
            <TextEditor id="spec" bind:value={spec} language="yaml" />
          </div>
          <details>
            <summary>Show/hide example</summary>
            <CodeBlock>
              {specExample}
            </CodeBlock>
          </details>
        {:else}
          <!-- GUI form edit mode -->
        {/if}
        <div class="buttons">
          <Button type="submit">Create Container</Button>
        </div>
        <div style="margin: 24px 0;">
          <ErrorLabel value={error} />
        </div>
      </form>
    </div>
  </Panel>
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

  .center-form {
    max-width: 640px;
    margin: 0 auto;
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
