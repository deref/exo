<script lang="ts">
  import Icon from './Icon.svelte';

  export let value: unknown = null;
</script>

{#if value != null}
  <div>
    <Icon glyph="Warning" />
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
