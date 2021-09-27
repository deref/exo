<script lang="ts">
  import Icon from './Icon.svelte';

  export let tooltip: string | undefined = undefined;
  export let glyph: string;
</script>

<button on:click title={tooltip}>
  <Icon {glyph} />
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
    color: var(--strong-color);
    background-color: var(--icon-button-hover-bg-color);
  }

  button:focus {
    background-color: var(--icon-button-focus-bg-color);
  }
</style>
