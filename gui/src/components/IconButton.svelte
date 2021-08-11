<script lang="ts">
  export let tooltip: string = '';
</script>

<button on:click title={tooltip !== '' ? tooltip : undefined}>
  <slot />
</button>

<style>
  button {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--grey-7-color);
    background-color: var(--icon-button-bg-color);
    border-radius: 4px;
    border: none;
    outline: none;
  }

  button :global(svg),
  button :global(svg *) {
    fill: currentColor;
  }

  button:hover {
    background-color: var(--icon-button-hover-bg-color);
    color: var(--grey-3-color);
  }

  button:focus {
    background-color: var(--icon-button-focus-bg-color);
  }
</style>
