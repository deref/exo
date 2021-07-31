<script lang="ts">
  export let message: string = '';

  const whitespaceRegex = /\s/;
  const linkRegex = /^https*:\/\/\S+\.\S+/;

  // Substring buckets of whitespace and non whitespace spans
  let buckets = [message[0]];

  for (let i = 1; i < message.length; i++) {
    if (
      whitespaceRegex.test(message[i]) !== whitespaceRegex.test(message[i - 1])
    ) {
      buckets.push(message[i]);
    } else {
      buckets[buckets.length - 1] += message[i];
    }
  }

  // Substring spans of links and non-links
  let spans = [buckets[0]];

  for (let i = 1; i < buckets.length; i++) {
    if (linkRegex.test(buckets[i]) !== linkRegex.test(buckets[i - 1])) {
      spans.push(buckets[i]);
    } else {
      spans[spans.length - 1] += buckets[i];
    }
  }
</script>

<span>
  {#each spans as span}
    {#if linkRegex.test(span)}
      <a href={span}>{span}</a>
    {:else}
      <span>{span}</span>
    {/if}
  {/each}
</span>

<style>
</style>
