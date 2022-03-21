<script lang="ts">
  import { location } from 'svelte-spa-router';
  import Panel from '../components/Panel.svelte';
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import CheckeredTable from '../components/CheckeredTable.svelte';
  import { mutation, query } from '../lib/graphql';

  export let params = { workspace: '' };
  const workspaceId = params.workspace;

  const q = query(
    `#graphql
    query ($workspaceId: String!) {
      workspace: workspaceById(id: $workspaceId) {
        id
        stack {
          vaults {
            id
            name
            url
            connected
            authenticated
          }
          environment {
            # TODO: Use a fragment.
            variables {
              name
              value
              source
            }
          }
        }
      }
    }`,
    {
      variables: {
        workspaceId,
      },
    },
  );

  // TODO: Should _remove_ vault from the stack, not necessarily forget it, but
  // need some global vaults GUI UX too before doing that.
  const forgetVaultMutation = mutation(
    `#graphql
    mutation ($id: String!) {
      forgetVault(ref: $id) {
        __typename
      }
    }`,
  );
  const forgetVault = async (id: string) => {
    await forgetVaultMutation({ variables: { id } });
  };

  const workspaceUrl = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const authEsv = async () => {
    const uri = new URL(window.location.href);
    uri.hash = '/auth-esv';
    uri.searchParams.set('returnTo', $location);
    window.location.href = uri.toString();
  };

  $: vaults = $q.data?.workspace?.stack?.vaults;
  $: variables = $q.data?.workspace?.stack?.environment?.variables;
</script>

<Layout loading={$q.loading} error={$q.error}>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <Panel title="Workspace Variables" backUrl={workspaceUrl}>
    <div class="vaults-title">
      <h2>Secrets Vaults</h2>
      <Button href={`${workspaceUrl}/add-vault`} small>+ Add vault</Button>
    </div>
    {#if vaults && vaults.length > 0}
      <CheckeredTable>
        <svelte:fragment slot="head">
          <th>Name</th>
          <th>URL</th>
        </svelte:fragment>
        {#each vaults as vault}
          <tr>
            <td>{vault.name}</td>
            <td>
              {#if vault.connected}
                <a href={vault.url} target="_blank">{vault.url}</a>
              {:else}
                {vault.url}
              {/if}
            </td>
            <td>
              {#if vault.connected}
                <Button href={`${vault.url}/create-secret`} small>
                  + New secret
                </Button>
              {:else if !vault.authenticated}
                <Button on:click={() => authEsv()} small>Authenticate</Button>
              {:else}
                Bad vault URL
              {/if}
              <Button
                on:click={async (_event) => {
                  await forgetVault(vault.id);
                  window.location.reload();
                }}
                small
              >
                x
              </Button>
            </td>
          </tr>
        {/each}
      </CheckeredTable>
    {:else if vaults}
      <div>No vaults linked to this workspace.</div>
    {/if}
    {#if variables && variables.length > 0}
      <hr />
      <h2>Variables</h2>
      <EnvironmentTable {variables} />
    {:else if variables}
      <div>Empty environment, no variables found.</div>
    {/if}
  </Panel>
</Layout>

<style>
  .vaults-title {
    display: flex;
    align-items: center;
    gap: 18px;
    margin-bottom: 24px;
  }

  .vaults-title h2 {
    margin: 0;
  }
</style>
