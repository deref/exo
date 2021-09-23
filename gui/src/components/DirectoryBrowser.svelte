<script lang="ts">
  import type { ReadDirResult } from '../lib/api';
  import HomeSVG from './mono/HomeSVG.svelte';
  import LeftUpSVG from './mono/LeftUpSVG.svelte';
  import Textbox from './Textbox.svelte';

  export let dir: ReadDirResult;
  export let homePath: string;
  export let handleClick: (path: string) => void;
  let search: string = '';
</script>

<div class="toolbar">
  {#if dir.parent !== null}
    <button
      title="Parent directory"
      on:click={() => handleClick(String(dir.parent?.path))}
    >
      <LeftUpSVG />
    </button>
  {:else}
    <button disabled> <LeftUpSVG /> </button>
  {/if}

  {#if dir.directory.path !== homePath}
    <button title="Home directory" on:click={() => handleClick(homePath)}>
      <HomeSVG />
    </button>
  {:else}
    <button disabled> <HomeSVG /> </button>
  {/if}

  <Textbox
    placeholder="Search..."
    bind:value={search}
    --input-width="100%"
    autofocus
  />
</div>

<div class="directories">
  {#each dir.entries
    .filter((x) => x.isDirectory)
    .filter((x) => x.name.slice(0, search.length) === search)
    .sort((x, y) => x.name.localeCompare(y.name)) as entry}
    <button title={entry.path} on:click={() => handleClick(entry.path)}>
      {entry.name}
    </button>
  {/each}
</div>

<style>
  .toolbar {
    display: grid;
    grid-template-columns: repeat(2, 48px) auto;
    gap: 12px;
  }

  .toolbar button {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  * :global(input) {
    height: 36px !important;
  }

  * :global(svg) {
    height: 16px;
  }

  * :global(svg *) {
    fill: currentColor !important;
  }

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
