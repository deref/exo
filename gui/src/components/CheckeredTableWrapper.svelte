<!-- 
  This component styles any <table> within the <slot /> with a checkered row-and-column tiger stripe
    as well as some shadowing and general style improvements for enhanced legibility.

  Please ensure your table has the structure:
  * table
    * thead (optional)
    * ...tr...
      * ...th...
    * tbody
    * ...tr...
      * ...td...
-->

<div>
  <slot />
</div>

<style>
  div :global(table) {
    border: none;
    border-collapse: collapse;
    border-radius: 4px;
    box-shadow: 0 0.33px 0 1px hsla(0, 0%, 100%, 0.15),
      0 6px 9px -4px rgba(0, 0, 0, 0.2), 0 0.4px 0 0.8px rgba(0, 0, 0, 0.1);
    overflow: hidden;
  }

  div :global(td),
  div :global(th) {
    text-align: left;
    padding: 6px 15px;
  }

  div :global(th:nth-child(2n)) {
    background: #f9f9f9;
  }

  div :global(tbody tr:nth-child(2n + 1)) {
    background: #eeeeee;
  }

  div :global(tbody tr:nth-child(2n)) {
    background: #ffffff;
  }

  div :global(tbody tr:nth-child(2n + 1) td:nth-child(2n)) {
    background: #e7e7e7;
  }

  div :global(tbody tr:nth-child(2n) td:nth-child(2n)) {
    background: #f9f9f9;
  }

  div :global(tr:first-child > *) {
    padding-top: 10px;
  }

  div :global(tr:last-child > *) {
    padding-bottom: 10px;
  }
</style>
