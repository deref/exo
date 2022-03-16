<script lang="ts" context="module">
  import type { Readable } from 'svelte/store';

  export type DirectoryStore = Readable<
    | {
        ready: true;
        directory: Directory;
        homePath: string;
      }
    | {
        ready: false;
        error: Error | null;
      }
  > & {
    setDirectory: (path: string) => void;
  };

  // TODO: Fragments

  export type Directory = {
    path: string;
    parentPath: null | string;
    children: File[];
  };

  export type File = {
    path: string;
    name: string;
    isDirectory: boolean;
  };
</script>

<script lang="ts">
  import { nonNull } from '../lib/util';
  import Icon from './Icon.svelte';
  import Textbox from './Textbox.svelte';
  import Spinner from '../components/Spinner.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';

  export let autofocus = false;
  export let store: DirectoryStore;

  let search: string = '';

  const goUp = () => {
    if (!$store.ready) {
      return;
    }
    const { parentPath } = $store.directory;
    if (parentPath) {
      store.setDirectory(parentPath);
    }
  };
  const goHome = () => {
    if (!$store.ready) {
      return;
    }
    store.setDirectory($store.homePath);
  };
</script>

{#if $store.ready}
  <div class="toolbar">
    <button
      title="Parent directory"
      on:click|preventDefault={goUp}
      disabled={$store.directory.parentPath === null}
    >
      <Icon glyph="LeftUp" />
    </button>

    <button
      title="Home directory"
      on:click|preventDefault={goHome}
      disabled={$store.directory.path === $store.homePath}
    >
      <Icon glyph="Home" />
    </button>

    <Textbox
      placeholder="Search..."
      bind:value={search}
      --input-width="100%"
      {autofocus}
    />
  </div>

  <div class="directories">
    {#each $store.directory.children
      .filter((x) => x.isDirectory)
      .filter((x) => x.name
          .toLowerCase()
          .startsWith(search.toLowerCase())) as file}
      <button
        title={file.path}
        on:click|preventDefault={() => store.setDirectory(file.path)}
      >
        {file.name}
      </button>
    {:else}
      <div class="no-match">
        {#if search}
          No subdirectories matching: <code>{search}</code>
        {:else}
          No subdirectories
        {/if}
      </div>
    {/each}
  </div>
{:else}
  <ErrorLabel value={$store.error} />
  <Spinner enabled={!$store.error} />
{/if}

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

  .no-match {
    padding: 4px;
    font-style: italic;
  }
  .no-match code {
    font-style: normal;
  }
</style>
