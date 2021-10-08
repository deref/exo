<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import EditAs from '../components/form/EditAs.svelte';
  import CodeBlock from '../components/CodeBlock.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import ShellEditor from '../components/ShellEditor.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import ArgumentsInput from '../components/ArgumentsInput.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import EnvironmentInput from '../components/EnvironmentInput.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import { parseScript, generateScript } from '../lib/process/script';
  import { setLogVisibility } from '../lib/logs/visible-logs';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;
  const workspaceNewComponentRoute = `/workspaces/${encodeURIComponent(
    workspaceId,
  )}/new-component`;

  let name: string = '';
  let script: string = '';
  let directory: string = '';
  let program: string = '';
  let environment: Record<string, string> = {};
  let args: string[] = [];

  const updateFields = () => {
    const parsed = parseScript(script);
    directory = parsed.spec.directory ?? '';
    program = parsed.spec.program ?? '';
    args = parsed.spec.arguments ?? [];
    environment = parsed.spec.environment ?? {};
    error = parsed.error;
  };

  const updateScript = () => {
    script = generateScript({
      directory,
      environment,
      program,
      arguments: args,
    });
  };

  let mode = 'script';
  const editorModes = [
    {
      id: 'script',
      name: 'Script',
      onActivate: updateScript,
    },
    {
      id: 'fields',
      name: 'Fields',
      onActivate: updateFields,
    },
  ];

  let error: Error | null = null;

  const codeExample = `# Export environment variables.
export DEBUG=true

# Set working directory.
cd /

# Specify command with arguments.
my-app --port 4000
`;
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <CenterFormPanel title="New Process" backRoute={workspaceNewComponentRoute}>
    <h1><Icon glyph="Layers" /> New Process</h1>
    <form
      on:submit|preventDefault={async () => {
        updateFields();
        try {
          const { id } = await workspace.createProcess(name, {
            directory,
            environment,
            program,
            arguments: args,
          });
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
      <div>
        <label for="name">Name:</label>
        <Textbox id="name" name="name" bind:value={name} />
      </div>

      <EditAs bind:mode {editorModes} />

      {#if mode === 'fields'}
        <div>
          <label for="program">Program:</label>
          <Textbox id="program" name="program" bind:value={program} />
        </div>
        <div>
          <label for="args">Arguments: (one per line)</label>
          <ArgumentsInput id="args" name="args" bind:value={args} />
        </div>
        <div>
          <label for="directory">Working Directory:</label>
          <Textbox id="directory" name="directory" bind:value={directory} />
        </div>
        <div>
          <label for="environment">Environment:</label>
          <EnvironmentInput
            id="environment"
            name="environment"
            bind:environment
          />
        </div>
        <SubmitButton>Create Process</SubmitButton>
      {:else}
        <div>
          <div class="script-editor">
            <label for="script">Script:</label>
            <ShellEditor id="script" bind:value={script} />
          </div>
          <details>
            <summary>Show/hide example</summary>
            <CodeBlock code={codeExample} />
          </details>
          <SubmitButton>Create Process</SubmitButton>
        </div>
      {/if}
      <ErrorLabel value={error} />
    </form>
  </CenterFormPanel>
</Layout>

<style>
  * :global(input),
  * :global(textarea),
  .script-editor {
    margin-bottom: 24px;
  }

  * :global(input) {
    width: 100%;
  }
</style>
