<script lang="ts">
  export let message: string = '';

  let words = message.split(' ');

  const linkRegex = /^https*:\/\/\S+\.\S+/;
</script>

<span>
  {#each words as word}
    {#if linkRegex.test(word)}
      <a href={word}>{word}</a>
    {:else}
      <span>{word}</span>
    {/if}
    {' '}
  {/each}
</span>

<style>
</style>
