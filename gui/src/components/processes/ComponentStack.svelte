<script lang="ts" context="module">
  import type { Component } from './ComponentCard.svelte';

  export type { Component };
</script>

<script lang="ts">
  import ComponentCard from './ComponentCard.svelte';

  export let components: Component[];
  export let setLogsVisible: (id: string, value: boolean) => void;
  export let setRun: (id: string, value: boolean) => void;
  export let dispose: (id: string) => void;
</script>

{#each components as component (component.id)}
  <ComponentCard
    {component}
    setLogsVisible={(value) => setLogsVisible(component.id, value)}
    setRun={(value) => setRun(component.id, value)}
    dispose={() => dispose(component.id)}
  />
{:else}
  <i>No components.</i>
{/each}
