<script lang="ts">
  import * as router from 'svelte-spa-router';

  export let href: string;
</script>

<a
  {href}
  on:click={(event) => {
    if (href.startsWith('#')) {
      router.push(href);
      event.preventDefault();
    }
  }}><slot /></a
>
