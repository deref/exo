<script lang="ts">
  import WarningSVG from './mono/WarningSVG.svelte';

  export let value: unknown = null;
</script>

{#if value !== null}
  <div>
    <WarningSVG />
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
    color: var(--error-color);
    display: flex;
    align-items: center;
  }

  span {
    display: block;
    padding-left: 8px;
  }
</style>
