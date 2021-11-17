<script lang="ts">
  import type { ReadDirResult } from '../lib/api';
  import Icon from './Icon.svelte';
  import Textbox from './Textbox.svelte';

  export let dir: ReadDirResult;
  export let homePath: string;
  export let handleClick: (path: string) => void;
  let search: string = '';
</script>

<div class="toolbar">
  <button
    title="Parent directory"
    on:click={() => handleClick(String(dir.parent?.path))}
    disabled={dir.parent === null}
  >
    <Icon glyph="LeftUp" />
  </button>

  <button
    title="Home directory"
    on:click={() => handleClick(homePath)}
    disabled={dir.directory.path === homePath}
  >
    <Icon glyph="Home" />
  </button>

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
    max-height: 20rem;
    overflow-y: auto;
    overflow-x: hidden;
    border-radius: 5px;
    box-shadow: var(--button-shadow);
  }

  .directories button {
    border-radius: 0;
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
    outline: none;
  }

  button:not(:last-child) {
    margin-bottom: 1px;
  }

  button:disabled {
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

  button:focus,
  button:focus-within {
    box-shadow: var(--shadow-focus) !important;
  }

  .directories button:focus,
  .directories button:focus-within {
    margin: 1px 2px;
    padding: 5px 8px;
  }
</style>
