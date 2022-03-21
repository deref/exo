<script lang="ts" context="module">
  export interface Variable {
    name: string;
    value: string;
    source?: string;
  }
</script>

<script lang="ts">
  import CheckeredTable from './CheckeredTable.svelte';

  export let variables: Variable[] = [];

  const hasSources = variables.some((v) => !!v.source);
</script>

<CheckeredTable>
  {#each variables as { name, value, source }}
    <tr>
      <td class="name">{name}</td>
      <td class="value"><code><pre>{value}</pre></code></td>
      {#if hasSources}
        <td>{source}</td>
      {/if}
    </tr>
  {/each}
</CheckeredTable>

<style>
  .name {
    font-size: 0.8em;
    font-weight: 450;
    color: var(--grey-5-color);
    min-width: 100px;
  }

  .value {
    min-width: 250px;
  }

  code {
    width: 100%;
    max-width: 600px;
    display: inline-block;
    overflow-x: auto;
    padding: 8px;
    margin: -10px;
  }
</style>
