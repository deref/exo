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
