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

  const isValidUrl = (url: string): boolean => {
    try {
      const uri = new URL(url);
      return (
        uri.hostname === 'localhost' || uri.hostname === 'secrets.deref.io'
      );
    } catch (_err: unknown) {
      return false;
    }
  };

  const authEsv = async () => {
    const uri = new URL(window.location.href);
    uri.hash = '/auth-esv';
    uri.searchParams.set('returnTo', `/workspaces/${workspace.id}/variables`);
    window.location.href = uri.toString();
  };
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <CenterFormPanel title="Add Vault" {backRoute}>
    <h1><Icon glyph="Lock" />Add Secrets Vault</h1>
    <form
      on:submit={async () => {
        if (!isValidUrl(vaultUrl)) {
          alert(
            ' The URL you entered is not a valid URL. Please enter a valid URL.',
          );
          return;
        }
        await workspace.addVault({
          url: vaultUrl,
        });
        const user = await kernel.getEsvUser(vaultUrl);
        if (!user) {
          await authEsv();
          return;
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
