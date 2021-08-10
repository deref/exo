<script lang="ts">
  export let title: string;
  export let active: string;
</script>

<button on:click class:active={active === title} {title}><slot /></button>

<style>
  button {
    position: relative;
    border: none;
    background: transparent;
    color: var(--grey-9-text-color);
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
    color: var(--grey-3-text-color);
  }

  button:hover {
    color: var(--grey-6-text-color);
  }

  button.active:hover {
    background: var(--nav-button-active-hover-bg-color);
    color: var(--grey-1-text-color);
  }
</style>
