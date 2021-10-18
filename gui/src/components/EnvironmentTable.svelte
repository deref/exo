<script lang="ts" context="module">
  export interface Variable {
    name: string;
    value: string;
    source?: string;
  }
</script>

<script lang="ts">
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';

  export let variables: Variable[] = [];

  const hasSources = variables.some((v) => !!v.source);

  const sorted = variables
    .sort((e1, e2) => (e1.source || '').localeCompare(e2.source || ''))
    .sort((e1, e2) => e1.name.localeCompare(e2.name));
</script>

<CheckeredTableWrapper>
  <tbody>
    <table>
      {#each sorted as { name, value, source }}
        <tr>
          <td class="label">{name}</td>
          <td><code><pre>{value}</pre></code></td>
          {#if hasSources}
            <td>{source}</td>
          {/if}
        </tr>
      {/each}
    </table>
  </tbody>
</CheckeredTableWrapper>

<style>
  .label {
    font-size: 0.8em;
    font-weight: 450;
    color: var(--grey-5-color);
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
