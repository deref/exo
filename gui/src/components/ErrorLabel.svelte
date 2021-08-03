<script lang="ts">
  import Warning from './mono/warning.svelte';

  export let value: unknown = null;
</script>

{#if value !== null}
  <div>
    <Warning />
    <span>
      {#if value instanceof Error}
        {value.message}
      {:else if value instanceof String}
        {value}
      {:else}
        Unknown Error
      {/if}
    </span>
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
