<script lang="ts">
  import { parseSpans, Span } from './parseSpans';

  export let message: string = '';

  const spans = parseSpans(message);
</script>

<span>
  {#each spans as span}
    {#if span.href != null}
      <a href={span.href}>{span.text}</a>
    {:else}
      <span>{span.text}</span>
    {/if}
  {/each}
</span>

<style>
</style>
