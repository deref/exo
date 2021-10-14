<script lang="ts">
  import type { VariableDescription } from 'src/lib/api';
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';

  export let variables: Record<string, VariableDescription> = {};
  const entries = Object.entries(variables)
    .map(([name, desc]) => ({ name, ...desc }))
    .sort((e1, e2) => e1.name.localeCompare(e2.name))
    .sort((e1, e2) => e1.source.localeCompare(e2.source));
</script>

<CheckeredTableWrapper>
  <tbody>
    <table>
      {#each entries as { name, value, source }}
        <tr>
          <td class="label">{name}</td>
          <td><code><pre>{value}</pre></code></td>
          <td>{source}</td>
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
