<script lang="ts">
  export let title: string | undefined = undefined;
  export let active: string | undefined = undefined;
</script>

<div>
  <button on:click class:active={title && active === title}>
    <slot />
  </button>
  {#if title}
    <span>{title}</span>
  {/if}
</div>

<style>
  div {
    position: relative;
  }

  div:not(:hover):not(:focus-within) span {
    display: none;
  }

  span {
    display: inline-block;
    position: absolute;
    left: calc(100% - 4px);
    top: calc(50% - 11px);
    font-size: 13px;
    font-weight: 400;
    color: var(--strong-color);
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    height: 22px;
    padding: 3px 6px;
    border-radius: 4px;
    white-space: nowrap;
  }
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
