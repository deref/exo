<script lang="ts">
  import { parseSpans } from './parseSpans';

  export let message: string = '';

  const spans = parseSpans(message);
</script>

<span>
  {#each spans as span}
    {#if span.type === 'link'}
      <a href={span.href}>{span.text}</a>
    {:else}
      <span>{span.text}</span>
    {/if}
  {/each}
</span>

<style>
</style>
