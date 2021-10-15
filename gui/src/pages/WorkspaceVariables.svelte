<script lang="ts">
  import * as router from 'svelte-spa-router';
  import Panel from '../components/Panel.svelte';
  import Button from '../components/Button.svelte';
  import Layout from '../components/Layout.svelte';
  import Spinner from '../components/Spinner.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';
  import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const makeRequests = () =>
    Promise.all([workspace.describeEnvironment(), workspace.describeVaults()]);
  let requests = makeRequests();

  const authEsv = async () => {
    const result = await api.kernel.authEsv();
    window.open(result.authUrl, '_blank')?.focus();

    // This alert is doing two jobs: informing the user of the auth code they
    // should see in Auth0 as well as providing an indication that the user has
    // finished authenticating when they dismiss the alert.
    alert(`You should see the following code in Auth0: ${result.authCode}`);
    requests = makeRequests();
  };
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <Panel title="Workspace Variables" backRoute={workspaceRoute}>
    {#await requests}
      <Spinner />
    {:then [variables, vaults]}
      <div class="vaults-title">
        <h2>Vaults</h2>
        <Button
          on:click={() => router.push(`${workspaceRoute}/add-vault`)}
          small
        >
          Add vault
        </Button>
      </div>
      {#if vaults.length}
        <CheckeredTableWrapper>
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>URL</th>
                <th />
              </tr>
            </thead>
            <tbody>
              {#each vaults as vault}
                <tr>
                  <td>{vault.name}</td>
                  <td>{vault.url}</td>
                  <td>
                    {#if vault.connected}
                      <Button href={`${vault.url}/create-secret`} small>
                        Create secret
                      </Button>
                    {:else if vault.needsAuth}
                      <Button on:click={authEsv} small>Authenticate</Button>
                    {:else}
                      Bad vault URL
                    {/if}
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </CheckeredTableWrapper>
      {:else}
        <div>No vaults linked to this workspace.</div>
      {/if}
      {#if Object.keys(variables).length === 0}
        <div>Empty Environment</div>
      {:else}
        <hr />
        <h2>Variables</h2>
        <EnvironmentTable {variables} />
      {/if}
    {/await}
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
