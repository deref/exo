<script lang="ts">
  export let type: string | undefined = undefined;

  export let small: boolean = false;

  export let inset: boolean = false;

  export let disabled = false;
</script>

<button {disabled} class:small class:inset on:click {type}>
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
    padding: 0.2em 0.4em;
    margin: 0 0.2em;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }

  .inset,
  .inset:hover,
  .inset:active {
    background: var(--button-inset-background);
    box-shadow: var(--button-inset-shadow);
  }
</style>
