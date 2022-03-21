<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import BytesLabel from '../components/BytesLabel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import CheckeredTable from '../components/CheckeredTable.svelte';
  import Sparkline from '../components/Sparkline.svelte';
  import { query } from '../lib/graphql';

  export let params = { workspace: '', component: '' };
  const workspaceId = params.workspace;
  const componentId = params.component;

  const q = query(
    `#graphql
    query ($componentId: String!) {
      component: componentById(id: $componentId) {
        id
        name
        asProcess {
          running
          cpuPercentage
        }
        environment {
          variables {
            name
            value
            source
          }
        }
      }
    }`,
    {
      variables: {
        componentId,
      },
      pollInterval: 1000,
    },
  );
  const component = $q.data?.component;
  const process = component?.asProcess;

  let cpuPercentages: number[] = [];
  $: {
    if (component?.cpuPercentage != null) {
      cpuPercentages = [
        ...cpuPercentages.slice(0, 99),
        component.cpuPercentage,
      ];
    }
  }
</script>

<Layout loading={$q.loading} error={$q.error}>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <Panel
    title={component ? component.name : 'Loading...'}
    backUrl={`/workspaces/${encodeURIComponent(workspaceId)}`}
  >
    {#if component}
      <CheckeredTable>
        <tr>
          <td class="label">Status</td>
          <td>{process.running ? 'Running' : 'Stopped'}</td>
          <td />
        </tr>
        {#if component.running}
          <tr>
            <td class="label">CPU</td>
            <td>
              {#if process.cpuPercent}
                {process.cpuPercent.toFixed(2)}%
              {/if}
            </td>
          </tr>
          <tr>
            <td class="label">Resident Memory</td>
            <td>
              {#if process.residentMemory}
                <BytesLabel value={process.residentMemory} />
              {/if}
            </td>
            <td />
          </tr>
          <tr>
            <td class="label">Started at</td>
            <td>
              {#if process.createTime}
                <span title={new Date(process.createTime).toISOString()}>
                  {new Date(process.createTime).toLocaleTimeString()}
                </span>
              {/if}
            </td>
            <td>
              {#if cpuPercentages.some((p) => p !== 0)}
                <Sparkline entries={cpuPercentages} interactive />
              {/if}
            </td>
          </tr>
          <tr>
            <td class="label">Local Ports</td>
            <td>{process.ports?.join(', ') ?? 'None'}</td>
            <td />
          </tr>
          <tr>
            <td class="label">Children</td>
            <td>
              {process.childrenExecutables?.join(', ') ?? 'None'}
            </td>
            <td />
          </tr>
        {/if}
      </CheckeredTable>
      <br />
      <h3>Environment</h3>
      <EnvironmentTable variables={component.environment.variables} />
    {/if}
  </Panel>
</Layout>

<style>
  .label {
    font-size: 0.8em;
    font-weight: 450;
    color: var(--grey-5-color);
  }
</style>
