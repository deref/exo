<script lang="ts">
  import Icon from './Icon.svelte';
  import type { IconGlyph } from './Icon.svelte';
  import * as router from 'svelte-spa-router';
  import { createEventDispatcher } from 'svelte';

  export let glyph: IconGlyph;
  export let href: string | undefined = undefined;
  export let danger: boolean | undefined = undefined;

  const dispatch = createEventDispatcher();

  const handleClick = (e: MouseEvent) => {
    if (href) {
      router.push(href);
      return;
    }
    dispatch('click', e);
  };
</script>

<button class:danger on:click|preventDefault={handleClick}>
  <Icon {glyph} />
  <slot />
</button>

<style>
  button :global(*) {
    fill: currentColor;
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
    outline: none;
  }

  button :global(svg) {
    height: 16px;
    margin-left: -8px;
  }

  button:hover,
  button:focus,
  button:focus-within {
    color: var(--strong-color);
    background: var(--grey-e-color);
  }

  .danger {
    color: var(--error-color-faded);
  }
  .danger:hover {
    color: var(--error-color);
  }
</style>
