<script lang="ts">
  import * as router from 'svelte-spa-router';
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import Textbox from '../components/Textbox.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import SubmitButton from '../components/form/SubmitButton.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const kernel = api.kernel;
  const backRoute = `/workspaces/${encodeURIComponent(workspaceId)}/variables`;

  // TODO: inject this.
  const esvUrl = 'https://secrets.deref.io/organizations';

  let vaultUrl: string = '';
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <CenterFormPanel title="Add Vault" {backRoute}>
    <h1><Icon glyph="Lock" /> Add Secrets Vault</h1>
    <form
      on:submit={async () => {
        await workspace.addVault({
          url: vaultUrl,
        });
        const user = await kernel.getEsvUser(vaultUrl);
        if (!user) {
          const result = await api.kernel.authEsv();
          window.open(result.authUrl, '_blank')?.focus();

          // This alert is doing two jobs: informing the user of the auth code they
          // should see in Auth0 as well as providing an indication that the user has
          // finished authenticating when they dismiss the alert.
          alert(
            `You should see the following code in Auth0: ${result.authCode}`,
          );
        }
        await router.push(
          `/workspaces/${encodeURIComponent(workspace.id)}/variables`,
        );
      }}
    >
      <p>Enter a URL for your vault:</p>
      <Textbox bind:value={vaultUrl} --input-width="100%" autofocus />
      <SubmitButton>Add Secrets Vault</SubmitButton>
    </form>

    <h2>Need a vault?</h2>
    <a href={esvUrl} target="_blank">Create one with Exo Secrets Vault</a>
  </CenterFormPanel>
</Layout>
