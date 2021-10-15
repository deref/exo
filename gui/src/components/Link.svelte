<script lang="ts">
  import * as router from 'svelte-spa-router';
  import { absoluteUrl } from '../lib/regex';

  export let href: string;
</script>

<a
  {href}
  on:click={(event) => {
    if (!absoluteUrl.test(href)) {
      router.push(href);
      event.preventDefault();
    }
  }}
>
  <slot />
</a>
