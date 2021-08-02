<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import { parseScript, generateScript } from '../lib/process/script';
  import { setLogVisibility } from '../lib/logs/visible-logs';
  import EnvironmentInput from '../components/EnvironmentInput.svelte';
  import ArgumentsInput from '../components/ArgumentsInput.svelte';
  import Textbox from '../components/Textbox.svelte';
  import Button from '../components/Button.svelte';
  import CodeBlock from '../components/CodeBlock.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import ShellEditor from '../components/ShellEditor.svelte';

  export let params = { workspace: '' };

  const workspace = api.workspace(params.workspace);
  const workspaceRoute = `/workspaces/${encodeURIComponent(params.workspace)}`;

  let name: string = '';

  let script: string = '';
  let structured = false;
  let error: Error | null = null;

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
</script>

<Layout showBackButton backButtonRoute={workspaceRoute}>
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
    <div class="edit-as">
      <span>Edit as:</span>
      <button
        class:selected={!structured}
        on:click|preventDefault={(e) => {
          if (!structured) {
            return;
          }
          structured = false;
          updateScript();
        }}
      >
        script
      </button>
      <button
        class:selected={structured}
        on:click|preventDefault={() => {
          if (structured) {
            return;
          }
          structured = true;
          updateFields();
        }}
      >
        fields
      </button>
    </div>
    {#if structured}
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
      <div class="buttons">
        <Button type="submit">Create Process</Button>
      </div>
    {:else}
      <div class="columns">
        <div>
          <div>
            <label for="script">Script:</label>
            <ShellEditor id="script" bind:value={script} />
          </div>
          <div class="buttons">
            <Button type="submit">Create Process</Button>
          </div>
        </div>
        <div>
          <label>Example:</label>
          <CodeBlock
            >{`# Export environment variables.
export DEBUG=true

# Set working directory.
cd /

# Specify command with arguments.
my-app --port 4000
`}
          </CodeBlock>
        </div>
      </div>
    {/if}
    <ErrorLabel value={error} />
  </form>
</Layout>

<style>
  form {
    padding: 40px;
    max-width: 980px;
    margin: 0 auto;
  }

  .columns {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 30px;
  }

  * :global(input) {
    width: 100%;
  }

  label {
    display: block;
    margin-top: 24px;
    margin-bottom: 8px;
  }

  .buttons {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin-top: 8px;
  }

  .edit-as {
    padding-top: 32px;
  }

  .edit-as span {
    margin-right: 8px;
  }

  .edit-as button {
    border: none;
    background: none;
    font-weight: 450;
    color: #777;
  }

  .edit-as button:hover {
    text-decoration: underline;
  }

  .edit-as .selected {
    color: #06e;
    font-weight: 450;
    text-decoration: underline;
  }
</style>
