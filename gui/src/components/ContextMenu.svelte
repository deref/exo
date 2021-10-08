<script lang="ts">
  export let title: string | undefined;
</script>

<nav>
  {#if title}<span>{title}</span>{/if}
  <slot />
</nav>

<style>
  nav {
    display: none;
    position: absolute;
    right: 0;
    background: var(--primary-bg-color);
    box-shadow: var(--dropdown-shadow);
    border-radius: 5px;
    padding: 4px 0;
    margin: -6px;
    z-index: 2;
  }

  nav > span {
    display: block;
    padding: 4px 12px;
    font-size: 0.8em;
    color: var(--grey-7-color);
  }
</style>
