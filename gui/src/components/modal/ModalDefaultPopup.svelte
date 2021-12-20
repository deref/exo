<script lang="ts">
  export let message = 'Missing message.';

  export let title: string | undefined;
</script>

{#if title}
  <h1>{title}</h1>
{/if}
<p>{message}</p>

<style>
  p {
    white-space: pre-wrap;
  }
</style>
