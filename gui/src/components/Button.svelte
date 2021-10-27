<script lang="ts">
  import * as router from 'svelte-spa-router';
  import { absoluteUrl } from '../lib/regex';

  export let type: string | undefined = undefined;
  export let href: string | undefined = undefined;

  export let small: boolean = false;

  export let danger: boolean = false;

  export let inset: boolean = false;

  export let disabled = false;

  const handleClick = (e: MouseEvent) => {
    if (href) {
      e.preventDefault();
      if (!absoluteUrl.test(href)) {
        // Handle internal routes.
        router.push(href);
        return;
      }
      // Handle external routes (open in new tab).
      window.open(href, '_blank')?.focus();
      return;
    }
  };
</script>

<button
  {disabled}
  class:small
  class:danger
  class:inset
  on:click={handleClick}
  on:click
  {type}
>
  <slot />
</button>

<style>
  button {
    border: none;
    border-radius: 5px;
    padding: 9px 18px;
    background: var(--button-background);
    box-shadow: var(--button-shadow);
  }

  .small {
    padding: 0.2em 0.4em;
    margin: 0 0.2em;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }

  .danger {
    color: var(--danger-button-color);
    background: var(--danger-button-background);
    box-shadow: var(--danger-button-shadow);
  }

  .danger:hover {
    background: var(--danger-button-hover-background);
    box-shadow: var(--danger-button-hover-shadow);
  }

  .danger:active {
    background: var(--danger-button-active-background);
    box-shadow: var(--danger-button-active-shadow);
  }

  .inset,
  .inset:hover,
  .inset:active {
    background: var(--button-inset-background);
    box-shadow: var(--button-inset-shadow);
  }
</style>
