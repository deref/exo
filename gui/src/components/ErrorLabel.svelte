<script lang="ts">
  import Warning from './mono/warning.svelte';

  export let value: Error | null = null;
</script>

{#if value}
  <div>
    <Warning />
    <span>{value.message}</span>
  </div>
{/if}

<style>
  div {
    display: flex;
    align-items: center;
    color: red;
  }

  span {
    display: block;
    padding-left: 8px;
  }
</style>
