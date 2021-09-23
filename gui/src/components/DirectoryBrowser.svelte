<script lang="ts">
  import type { ReadDirResult } from '../lib/api';

  export let dir: ReadDirResult;
  export let handleClick: (path: string) => void;
</script>

{#if dir.parent !== null}
  <button on:click={() => handleClick(String(dir.parent?.path))}> .. </button>
{:else}
  <button disabled> .. </button>
{/if}
<div class="directories">
  {#each dir.entries
    .filter((x) => x.isDirectory)
    .sort((x, y) => x.name.localeCompare(y.name)) as entry}
    <button on:click={() => handleClick(entry.path)}>{entry.name}</button>
  {/each}
</div>

<style>
  .directories {
    margin: 12px 0;
  }

  .directories button {
    border-radius: 0;
  }

  .directories button:first-of-type {
    border-top-left-radius: 5px;
    border-top-right-radius: 5px;
  }

  .directories button:last-of-type {
    border-bottom-left-radius: 5px;
    border-bottom-right-radius: 5px;
  }

  button {
    background: var(--primary-bg-color);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 5px;
    padding: 5px 10px;
    position: relative;
    display: grid;
    width: 100%;
    grid-template-columns: max-content 2fr;
    align-items: center;
    gap: 12px;
    margin-top: 1px;
  }

  button:disabled {
    cursor: default;
    opacity: 0.5;
  }

  button:not(:disabled):hover {
    background: var(--grey-e-color);
    box-shadow: var(--button-hover-shadow);
  }

  button:not(:disabled):active {
    background: var(--grey-c-color);
    box-shadow: var(--button-active-shadow);
  }
</style>
