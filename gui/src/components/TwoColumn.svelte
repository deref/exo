<div>
  <slot />
  <slot />
</div>

<style>
  div {
    height: 100%;
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 1px;
  }

  div > :global(:first-child) {
    --panel-bg-color: var(--secondary-bg-color);
  }
</style>
