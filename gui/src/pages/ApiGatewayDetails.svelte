<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { api, ApiGatewayDescription } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';
  import { fetchApiGateways, apiGateways } from '../lib/process/store';
  import type { RequestLifecycle } from '../lib/api';
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';

  export let params = { workspace: '', component: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const componentId = params.component;
  let apiGatewayList: RequestLifecycle<ApiGatewayDescription[]> = {
    stage: 'pending',
  };
  const unsubscribeApiGateways = apiGateways.subscribe((ag) => {
    apiGatewayList = ag;
  });

  let refreshInterval: ReturnType<typeof setInterval>;
  let apiGateway: ApiGatewayDescription | null = null;

  interface Endpoint {
    name: string;
    url: string;
  }
  let endpoints: Endpoint[] = [];

  onMount(() => {
    fetchApiGateways(workspace);
    refreshInterval = setInterval(() => {
      workspace
        .describeProcesses()
        .then((procs) => {
          if (apiGateway) {
            endpoints = [];
            for (const p of procs) {
              endpoints.push({
                name: p.name,
                url: `${p.name}.exo.localhost:${apiGateway.apiPort}`,
              });
            }
          }
        })
        .catch((err) => {
          console.error(err);
        });
      if (apiGatewayList.stage === 'success') {
        apiGateway = apiGatewayList.data.filter((p) => p.id === componentId)[0];
      }
    }, 1000);
  });

  onDestroy(() => {
    clearInterval(refreshInterval);
    unsubscribeApiGateways();
  });
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  {#if apiGateway}
    <Panel title={apiGateway.name} backRoute={workspaceRoute}>
      <h1>{apiGateway.name} API Gateway</h1>
      {#if apiGateway.running}
        <CheckeredTableWrapper>
          <tbody>
            <table>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>URL</th>
                </tr>
              </thead>
              {#each endpoints as { name, url }}
                <tr>
                  <td>{name}</td>
                  <td><a href={'//' + url} target="_blank">{url}</a></td>
                </tr>
              {:else}
                <tr><td colspan="2"> Loading... </td></tr>
              {/each}
            </table>
          </tbody>
        </CheckeredTableWrapper>
        <p>
          <a href={`http://localhost:${apiGateway.webPort}`} target="_blank"
            >Web interface:</a
          >
        </p>
        <iframe
          src={`http://localhost:${apiGateway.webPort}/#/flows`}
          title="Connections"
        />
        <br />
      {:else}
        <span>API gateway is not running</span>
      {/if}
    </Panel>
  {:else}
    <Panel title="Loading..." backRoute={workspaceRoute} />
  {/if}
</Layout>

<style>
  iframe {
    margin-top: 1em;
    border: none;
    width: 100%;
    min-height: 400px;
  }

  p {
    margin-top: 3em;
    margin-bottom: 0;
  }
</style>
