<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Code from '../components/Code.svelte';
  import { api } from '../lib/api';

  const workspaces = api.kernel.describeWorkspaces();
</script>

<Layout>
  <section>
    Use <Code>exo gui</Code> in your terminal.
    <!-- TODO: Use this data below: !-->
    <div style="display: none">
      {#await workspaces}
        loading workspaces...
      {:then workspaces}
        <ul>
          {#each workspaces as workspace}
            <li>{workspace.id} - {workspace.root}</li>
          {/each}
        </ul>
      {:catch error}
        <p style="color: red">{error.message}</p>
      {/await}
    </div>
  </section>
</Layout>

<style>
  section {
    padding: 30px;
  }
</style>
