<script lang="ts">
  import NavbarButton from './NavbarButton.svelte';
  import * as router from 'svelte-spa-router';

  export let title: string | undefined;
  export let active: string | undefined;
  export let href: string;
</script>

<NavbarButton
  {title}
  {active}
  on:click={() => {
    router.push(href);
  }}
>
  <slot />
</NavbarButton>
