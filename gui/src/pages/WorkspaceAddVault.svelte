<script lang="ts">
  import * as router from 'svelte-spa-router';
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  let vaultUrl: string = '';
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <Panel title="Workspace Variables" backRoute={workspaceRoute}>
    <form
      on:submit={async () => {
        await workspace.addVault({
          name: 'exo-vault',
          url: vaultUrl,
        });
        await router.push(
          `/workspaces/${encodeURIComponent(workspace.id)}/variables`,
        );
      }}
    >
      <p>Enter a URL for your vault:</p>
      <Textbox bind:value={vaultUrl} --input-width="100%" autofocus />
      <button type="submit">Add</button>
    </form>
  </Panel>
</Layout>
