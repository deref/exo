<script lang="ts">
  import Button from '../Button.svelte';

  export let procfileExport: any;
  export let clickHandler: any;
</script>

{#if procfileExport}
  <p>
    Your Procfile is not up to date.
    <Button small on:click={clickHandler}>Export?</Button>
  </p>
{/if}

<style>
  p {
    margin: 24px 0;
  }
</style>
