<script lang="ts">
  import type { Optional } from '../lib/util';

  // In range 0 to 1.
  export let value: Optional<number>;

  // Number of places after the decimal.
  export let precision = 0;
</script>

{#if value != null}
  <span>{value.toFixed(precision)}%</span>
{/if}
