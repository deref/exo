<script lang="ts">
  import Tooltip from '../Tooltip.svelte';
  import type { ShortcutParams } from '../../lib/actions/shortcut';
  import { shortcut } from '../../lib/actions/shortcut';
  import * as router from 'svelte-spa-router';
  import { createEventDispatcher } from 'svelte';

  export let title: string | undefined = undefined;
  export let active: string | undefined = undefined;
  export let href: string | undefined = undefined;
  export let shortcutParams: ShortcutParams | undefined = undefined;

  const dispatch = createEventDispatcher();

  const handleClick = (e: MouseEvent) => {
    if (href) {
      router.push(href);
      return;
    }
    dispatch('click', e);
  };
</script>

<div>
  {#if shortcutParams !== undefined}
    <button
      use:shortcut={shortcutParams}
      on:click|preventDefault={handleClick}
      class:active={title && active === title}
    >
      <slot />
    </button>
  {:else}
    <button
      on:click|preventDefault={handleClick}
      class:active={title && active === title}
    >
      <slot />
    </button>
  {/if}
  {#if title}
    <div class="tooltip"><Tooltip>{title}</Tooltip></div>
  {/if}
</div>

<style>
  div {
    position: relative;
  }

  div:not(:hover):not(:focus-within) .tooltip {
    display: none;
  }

  .tooltip {
    position: absolute;
    left: calc(100% - 4px);
    top: calc(50% - 11px);
  }
  button {
    position: relative;
    border: none;
    background: transparent;
    color: var(--grey-9-color);
    height: 48px;
    width: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  button > :global(svg),
  button > :global(svg *) {
    fill: currentColor;
  }

  button.active {
    background: var(--nav-button-active-bg-color);
    color: var(--grey-3-color);
  }

  button:hover {
    color: var(--grey-6-color);
  }

  button.active:hover {
    background: var(--nav-button-active-hover-bg-color);
    color: var(--grey-1-color);
  }
</style>
