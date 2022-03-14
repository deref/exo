<script lang="ts">
  import ComponentCard from './ComponentCard.svelte';

  export let components: {
    id: string;
    name: string;
    reconciling: boolean;
    running: boolean;
    logsVisible: boolean;
  }[];
  export let setRun: (id: string, value: boolean) => void;
  export let showPath: (id: string) => string;
  export let editPath: (id: string) => string;
  export let dispose: (id: string) => void;
</script>

{#each components as { id, name, running, reconciling, logsVisible } (id)}
  <ComponentCard
    {name}
    {running}
    {reconciling}
    {logsVisible}
    setLogsVisible={(value) => {
      setRun(id, value);
    }}
    showPath={showPath(id)}
    editPath={editPath(id)}
    setRun={(value) => {
      setRun(id, value);
    }}
    dispose={() => dispose(id)}
  />
{:else}
  <i>No components.</i>
{/each}
