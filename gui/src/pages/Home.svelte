<script lang="ts">
  import Code from '../components/Code.svelte';
  import Layout from '../components/Layout.svelte';
  import WorkspaceList from '../components/WorkspaceList.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { query } from '../lib/graphql';

  const q = query(`#graphql
    {
      workspaces: allWorkspaces {
        id
        root
        displayName
      }
    }
  `);

  $: workspaces = $q.data?.workspaces;
</script>

<Layout loading={$q.loading} error={$q.error}>
  <CenterFormPanel title="Workspaces">
    <h1>Workspaces</h1>
    <div>
      {#if workspaces}
        <WorkspaceList {workspaces} />
      {/if}
    </div>
    <hr />
    <div>
      Use <Code>exo gui</Code> in your terminal to launch into the current directory's
      workspace.
    </div>
  </CenterFormPanel>
</Layout>
