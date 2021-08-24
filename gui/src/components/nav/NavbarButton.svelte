<script lang="ts">
  export let title: string | undefined = undefined;
  export let active: string | undefined = undefined;
</script>

<button on:click class:active={title && active === title} {title}>
  <slot />
</button>

<style>
  button {
    position: relative;
    border: none;
    background: transparent;
    color: var(--grey-9-color);
    height: 48px;
    width: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  button > :global(svg),
  button > :global(svg *) {
    fill: currentColor;
  }

  button.active {
    background: var(--nav-button-active-bg-color);
    color: var(--grey-3-color);
  }

  button:hover {
    color: var(--grey-6-color);
  }

  button.active:hover {
    background: var(--nav-button-active-hover-bg-color);
    color: var(--grey-1-color);
  }
</style>
