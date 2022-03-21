<!-- 
  Table with a checkered row-and-column tiger stripe, shadowing, and legibility
  improvements.
  
  Slot "head" requires a svelte:fragment of th elements.
  Default slot requires tr elements containing td elements.
-->

<table>
  {#if $$slots.head}
    <thead>
      <tr>
        <slot name="head" />
      </tr>
    </thead>
  {/if}
  <tbody>
    <slot />
  </tbody>
</table>

<style>
  table {
    border: none;
    border-collapse: collapse;
    border-radius: 4px;
    box-shadow: var(--heavy-3d-box-shadow);
    overflow: hidden;
  }

  table :global(th) {
    font-size: 0.85em;
    font-weight: 450;
    color: var(--grey-5-color);
  }

  table :global(td),
  table :global(th) {
    text-align: left;
    padding: 6px 15px;
  }

  table :global(th:nth-child(2n)) {
    background: var(--grey-f9-color);
  }

  table :global(tbody tr:nth-child(2n + 1)) {
    background: var(--grey-e-color);
  }

  table :global(tbody tr:nth-child(2n)) {
    background: var(--primary-bg-color);
  }

  table :global(tbody tr:nth-child(2n + 1) td:nth-child(2n)) {
    background: var(--grey-e7-color);
  }

  table :global(tbody tr:nth-child(2n) td:nth-child(2n)) {
    background: var(--grey-f9-color);
  }

  table :global(tr:first-child > *) {
    padding-top: 10px;
  }

  table :global(tr:last-child > *) {
    padding-bottom: 10px;
  }
</style>
