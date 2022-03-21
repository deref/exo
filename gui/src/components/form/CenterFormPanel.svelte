<script lang="ts">
  import Panel from '../Panel.svelte';

  // This component is an opinionated wrapper for creating
  // centered, vertically scrolling form layouts in exo.
  // It is to be used for most creation/edit forms and
  // intermediate step pages in create/edit flows.

  export let title: string = '';
  export let backUrl: string = '';
</script>

<Panel {title} {backUrl} --panel-padding="2rem" --panel-overflow-y="scroll">
  <div slot="actions">
    <slot name="actions" />
  </div>
  <div>
    <slot />
  </div>
</Panel>

<style>
  div {
    max-width: 640px;
    margin: 0 auto;
  }
</style>
