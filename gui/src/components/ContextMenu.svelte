<script lang="ts">
  import Icon from './Icon.svelte';
  import type { IconGlyph } from './Icon.svelte';

  export let title: string | undefined;

  type Action = { name: string; glyph: IconGlyph; function: any };

  export let actions: Action[];
</script>

<nav>
  {#if title}<span>{title}</span>{/if}
  {#each actions as action}
    <button on:click={action.function}>
      <Icon glyph={action.glyph} />
      {action.name}
    </button>
  {/each}
</nav>

<style>
  nav {
    display: none;
    position: absolute;
    right: 0;
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    border-radius: 5px;
    padding: 4px 0;
    margin: -6px;
    z-index: 2;
  }

  nav > span {
    display: block;
    padding: 4px 12px;
    font-size: 0.8em;
    color: var(--grey-7-color);
  }

  button {
    background: none;
    border: none;
    display: flex;
    align-items: center;
    font-size: 0.9em;
    gap: 4px;
    border-radius: 2px;
    padding: 6px 18px;
    width: 100%;
    white-space: nowrap;
    color: var(--grey-5-color);
  }

  button :global(*) {
    fill: currentColor;
  }

  button :global(svg) {
    height: 16px;
    margin-left: -8px;
  }

  button:hover {
    color: var(--strong-color);
    background: var(--grey-e-color);
  }
</style>
