<!-- 
  This component styles any <table> within the <slot /> with a checkered row-and-column tiger stripe
    as well as some shadowing and general style improvements for enhanced legibility.

  Please ensure your table has the structure:
  table
      thead (optional)
          ...tr...
              ...th...
      tbody
          ...tr...
              ...td...
-->

<div>
  <slot />
</div>

<style>
  div :global(table) {
    border: none;
    border-collapse: collapse;
    border-radius: 4px;
    box-shadow: var(--heavy-3d-box-shadow);
    overflow: hidden;
  }

  div :global(th) {
    font-size: 0.85em;
    font-weight: 450;
    color: var(--grey-5-color);
  }

  div :global(td),
  div :global(th) {
    text-align: left;
    padding: 6px 15px;
  }

  div :global(th:nth-child(2n)) {
    background: var(--grey-f9-color);
  }

  div :global(tbody tr:nth-child(2n + 1)) {
    background: var(--grey-e-color);
  }

  div :global(tbody tr:nth-child(2n)) {
    background: var(--primary-bg-color);
  }

  div :global(tbody tr:nth-child(2n + 1) td:nth-child(2n)) {
    background: var(--grey-e7-color);
  }

  div :global(tbody tr:nth-child(2n) td:nth-child(2n)) {
    background: var(--grey-f9-color);
  }

  div :global(tr:first-child > *) {
    padding-top: 10px;
  }

  div :global(tr:last-child > *) {
    padding-bottom: 10px;
  }
</style>
