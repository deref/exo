<script lang="ts">
  import { api } from '../lib/api';
  import * as router from 'svelte-spa-router';
  import { parseScript, generateScript } from '../lib/process/script';
  import EnvironmentInput from '../components/EnvironmentInput.svelte';
  import ArgumentsInput from '../components/ArgumentsInput.svelte';
  
  export let params = { workspace: '' };
  
  const workspace = api.workspace(params.workspace);
  
  let name: string = '';

  let script: string = '';
  let structured = false;
  let error: string | null = null;
  
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

<div class="layout">
  <form on:submit|preventDefault={async () => {
    updateFields();
    await workspace.createProcess(name, {
      directory,
      environment,
      program,
      arguments: args,
    });
    router.push(`/workspaces/${encodeURIComponent(params.workspace)}`);
  }}>
    <label>
      Name:
      <input bind:value={name} name="name"/>
    </label>
    {#if structured}
      <button on:click={() => {
        structured = false;
        updateScript();
      }}>
        Edit Script
      </button>
      <label>
        Program:
        <input name="program" bind:value={program}/>
      </label>
      <label for="args">
        Arguments: (one per line)
        <ArgumentsInput name="args" bind:value={args}/>
      </label>
      <label>
        Working Directory:
        <input name="directory" bind:value={directory}/>
      </label>
      <label for="environment">
        Environment:
        <EnvironmentInput name="environment" bind:environment={environment}/>
      </label>
    {:else}
      <button on:click={() => {
        structured = true;
        updateFields();
      }}>
        Edit Fields
      </button>
      <label>
        Script:
        <textarea name="script" bind:value={script}/>
      </label>
      <div>
        Example:
        <pre>
# Export environment variables.
export DEBUG=true

# Set working directory.
cd /

# Specify command with arguments.
my-app --port 4000
        </pre>
      </div>
    {/if}
    <button type="submit">Create Process</button>
  </form>
</div>

<style>

  .layout {
    display: grid;
    grid-template-columns: 2fr 3fr;
    gap: 30px;
    margin: 30px;
  }
  
  label {
    display: block;
  }
  
  [name=script] {
    width: 400px;
    height: 100px;
  }
  
  button {
    margin: 16px;
  }

</style>
