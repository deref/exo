<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import { api, isClientError } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import { parseScript, generateScript } from '../lib/process/script';
  import EnvironmentInput from '../components/EnvironmentInput.svelte';
  import ArgumentsInput from '../components/ArgumentsInput.svelte';
  import Textbox from '../components/Textbox.svelte';
  import Textarea from '../components/Textarea.svelte';
  import Button from '../components/Button.svelte';
  import CodeBlock from '../components/CodeBlock.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  
  export let params = { workspace: '' };
  
  const workspace = api.workspace(params.workspace);
  
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
    directory = parsed.spec.directory;
    program = parsed.spec.program;
    args = parsed.spec.arguments;
    environment = parsed.spec.environment;
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

<Layout>
  <form on:submit|preventDefault={async () => {
    updateFields();
    try {
      await workspace.createProcess(name, {
        directory,
        environment,
        program,
        arguments: args,
      });
      router.push(`/workspaces/${encodeURIComponent(params.workspace)}`);
    } catch (ex) {
      if (!isClientError(ex)) {
        throw ex;
      }
      error = ex;
    }
  }}>
    <div>
    <label for="name">Name:</label>
      <Textbox id="name" name="name" bind:value={name} />
    </div>
    <div class="edit-as">
      Edit as:
      <button class:selected={!structured} on:click|preventDefault={(e) => {
        if (!structured) {
          return;
        }
        structured = false;
        updateScript();
      }}>
        script
      </button>
      <button class:selected={structured} on:click|preventDefault={() => {
        if (structured) {
          return;
        }
        structured = true;
        updateFields();
      }}>
        fields
      </button>
    </div>
    {#if structured}
      <div>
        <label for="program">Program:</label>
        <Textbox id="program" name="program" bind:value={program}/>
      </div>
      <div>
        <label for="args">Arguments: (one per line)</label>
        <ArgumentsInput id="args" name="args" bind:value={args}/>
      </div>
      <div>
        <label for="directory">Working Directory:</label>
        <Textbox id="directory" name="directory" bind:value={directory}/>
      </div>
      <div>
        <label>Environment:</label>
        <EnvironmentInput name="environment" bind:environment={environment}/>
      </div>
      <div class="buttons">
        <Button type="submit">Create Process</Button>
      </div>
    {:else}
      <div class="columns">
        <div>
          <div>
            <label for="script">Script:</label>
            <Textarea id="script" name="script" bind:value={script}/>
          </div>
          <div class="buttons">
            <Button type="submit">Create Process</Button>
          </div>
        </div>
        <div>
          <label>Example:</label>
          <CodeBlock>
# Export environment variables.
export DEBUG=true

# Set working directory.
cd /

# Specify command with arguments.
my-app --port 4000
          </CodeBlock>
        </div>
      </div>
    {/if}
    <ErrorLabel value={error}/>
  </form>
</Layout>

<style>

form {
  padding: 40px;
  max-width: 1000px;
  margin: 0 auto;
}

.columns {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
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

.edit-as button {
  border: none;
  background: none;
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
