<script lang="ts">
  export let type: string | undefined = undefined;

  // TODO: Something better than this.
  export let size: 'small' | 'medium' = 'medium';

  export let disabled = false;
</script>

<button {disabled} class={size} on:click {type}>
  <slot />
</button>

<style>
  button {
    border: none;
    border-radius: 5px;
    padding: 12px 18px;
    background: var(--button-background);
    box-shadow: var(--button-shadow);
  }

  .small {
    padding: 3px 10px;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }
</style>
