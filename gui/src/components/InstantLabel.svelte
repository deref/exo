<script lang="ts">
  import type { Optional } from '../lib/util';

  export let value: Optional<Instant>;
</script>

{#if value}
  <span title={new Date(value).toISOString()}>
    {new Date(value).toLocaleTimeString()}
  </span>
{/if}
