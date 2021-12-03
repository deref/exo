<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import { api, ApiGatewayDescription } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';
  import { fetchApiGateways, apiGateways } from '../lib/process/store';
  import type { RequestLifecycle } from '../lib/api';

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
                url: `${p.name}.exo.localhost:${apiGateway.webPort}`,
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
      <p>{JSON.stringify(apiGateway)}</p>
      <p>{JSON.stringify(endpoints)}</p>
      {#if apiGateway.running}
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
</style>
