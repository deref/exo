<script lang="ts">
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';

  export let variables: Record<string, string> = {};
</script>

<CheckeredTableWrapper>
  <tbody>
    <table>
      {#each Object.keys(variables ?? []).sort() as name (name)}
        <tr>
          <td class="label">{name}</td>
          <td><code><pre>{variables[name]}</pre></code></td>
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
