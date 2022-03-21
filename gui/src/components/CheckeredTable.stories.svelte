<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import CheckeredTable from './CheckeredTable.svelte';
</script>

<Meta title="CheckeredTable" component={CheckeredTable} />

<Template let:args>
  <CheckeredTable>
    <svelte:fragment slot="head">
      {#if args.header}
        <th>Left</th>
        <th>Right</th>
      {/if}
    </svelte:fragment>
    <tr>
      <td>One</td>
      <td>Two</td>
    </tr>
    <tr>
      <td>Three</td>
      <td>Four</td>
    </tr>
    <tr>
      <td>Five</td>
      <td>Six</td>
    </tr>
  </CheckeredTable>
</Template>

<Story name="Headered" args={{ header: true }} />
<Story name="Headless" args={{ header: false }} />
