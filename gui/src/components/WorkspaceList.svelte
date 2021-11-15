<script lang="ts">
  import * as router from 'svelte-spa-router';
  import type { WorkspaceDescription } from '../lib/api';

  export let workspaces: WorkspaceDescription[];
</script>

<button on:click={() => router.push('#/new-project')}> Add new project </button>

<div>
  {#each workspaces as workspace}
    <button
      on:click={() => {
        router.push(`/workspaces/${encodeURIComponent(workspace.id)}`);
      }}
      title={workspace.root}
    >
      <b>{workspace.displayName}</b>
      <span>{workspace.root}</span>
    </button>
  {/each}
</div>

<style>
  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 16px 24px;
    position: relative;
    display: grid;
    width: 100%;
    grid-template-columns: max-content 2fr;
    align-items: center;
    text-align: left;
    gap: 24px;
    margin-top: 12px;
  }

  button > b {
    font-weight: 550;
  }

  button > span {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    text-align: right;
    font-size: 0.875em;
    color: var(--grey-5-color);
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
