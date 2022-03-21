<script lang="ts">
  import { location } from 'svelte-spa-router';
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
    const uri = new URL(window.location.href);
    uri.hash = '/auth-esv';
    uri.searchParams.set('returnTo', $location);
    window.location.href = uri.toString();
  };
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Variables" slot="navbar" />
  <Panel title="Workspace Variables" backUrl={workspaceRoute}>
    {#await requests}
      <Spinner />
    {:then [variables, vaults]}
      <div class="vaults-title">
        <h2>Secrets Vaults</h2>
        <Button href={`${workspaceRoute}/add-vault`} small>+ Add vault</Button>
      </div>
      {#if vaults.length > 0}
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
                    {:else if vault.needsAuth}
                      <Button on:click={() => authEsv()} small
                        >Authenticate</Button
                      >
                    {:else}
                      Bad vault URL
                    {/if}
                    <Button
                      on:click={async (_event) => {
                        await workspace.removeVault({
                          url: vault.url,
                        });
                        window.location.reload();
                      }}
                      small
                    >
                      x
                    </Button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </CheckeredTableWrapper>
      {:else}
        <div>No vaults linked to this workspace.</div>
      {/if}
      {#if Object.keys(variables).length > 0}
        <hr />
        <h2>Variables</h2>
        <EnvironmentTable
          variables={Object.entries(variables).map(([name, description]) => ({
            name,
            ...description,
          }))}
        />
      {:else}
        <div>Empty environment, no variables found.</div>
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
