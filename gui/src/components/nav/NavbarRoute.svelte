<script lang="ts">
  import NavbarButton from './NavbarButton.svelte';
  import * as router from 'svelte-spa-router';
  import type { ShortcutParams } from '../../lib/actions/shortcut';

  export let title: string | undefined = undefined;
  export let active: string | undefined = undefined;
  export let href: string;

  export let shortcutParams: ShortcutParams | undefined = undefined;
</script>

<NavbarButton
  {title}
  {active}
  {shortcutParams}
  on:click={() => {
    router.push(href);
  }}
>
  <slot />
</NavbarButton>
