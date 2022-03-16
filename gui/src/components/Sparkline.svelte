<script lang="ts">
  import sparkline from '@fnando/sparkline';
  import { afterUpdate } from 'svelte';

  export let entries: number[];
  export let interactive = false;

  let svg: SVGSVGElement;

  afterUpdate(() => {
    sparkline(svg, entries, {
      interactive,
    });
  });
</script>

<svg
  bind:this={svg}
  class="sparkline"
  width="100"
  height="30"
  stroke-width="3"
/>

<style>
  .sparkline {
    margin: -6px -15px;
    stroke: var(--sparkline-stroke);
    fill: var(--sparkline-fill);
  }
</style>
