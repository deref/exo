<script lang="ts">
  import * as router from 'svelte-spa-router';
  import type { WorkspaceDescription } from '../lib/api';

  export let workspaces: WorkspaceDescription[];
</script>

<button on:click={() => router.push('#/new-workspace')}>
  Add new project
</button>

{#each workspaces as workspace}
  <button
    on:click={() => {
      router.push(`/workspaces/${encodeURIComponent(workspace.id)}`);
    }}
  >
    <b>{workspace.id}</b>
    <span title={workspace.root}>{workspace.root}</span>
  </button>
{/each}

<style>
  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 16px 24px;
    position: relative;
    display: grid;
    grid-template-columns: max-content 2fr;
    align-items: center;
    gap: 12px;
    margin-top: 12px;
  }

  button > b {
    max-width: 6em;
  }

  button > * {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }
</style>
