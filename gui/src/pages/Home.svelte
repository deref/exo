<script lang="ts">
  import Code from '../components/Code.svelte';
  import Layout from '../components/Layout.svelte';
  import ErrorLabel from '../components/ErrorLabel.svelte';
  import WorkspaceList from '../components/WorkspaceList.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { query } from '../lib/graphql';
  import { nonNull } from '../lib/util';

  const q = query(`#graphql
    {
      workspaces: allWorkspaces {
        id
        root
        displayName
      }
    }
  `);
</script>

<Layout>
  <CenterFormPanel title="Workspaces">
    <h1>Workspaces</h1>
    <div>
      {#if $q.loading}
        loading workspaces...
      {:else if $q.error}
        <ErrorLabel value={$q.error} />
      {:else}
        <WorkspaceList workspaces={nonNull($q.data).workspaces} />
      {/if}
    </div>
    <hr />
    <div>
      Use <Code>exo gui</Code> in your terminal to launch into the current directory's
      workspace.
    </div>
  </CenterFormPanel>
</Layout>
