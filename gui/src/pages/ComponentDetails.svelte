<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import BytesLabel from '../components/BytesLabel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import CheckeredTable from '../components/CheckeredTable.svelte';
  import Sparkline from '../components/Sparkline.svelte';
  import PercentLabel from '../components/PercentLabel.svelte';
  import InstantLabel from '../components/InstantLabel.svelte';
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
        running
        asProcess {
          cpuPercent
          residentBytes
          started
          ports
          environment {
            ...EnvironmentVariables
          }
        }
        environment {
          ...EnvironmentVariables
        }
      }
    }
    
    # TODO: Lift to EnvironmentTable.
    fragment EnvironmentVariables on Environment {
      variables {
        name
        value
        source
      }
    }
    `,
    {
      variables: {
        componentId,
      },
      pollInterval: 1000,
    },
  );
  const component = $q.data?.component;
  const process = component?.asProcess;
  // TODO: Visually distinguish observed environment vs specified environment.
  const environment = process?.environment ?? component?.environment;

  let cpuPercents: number[] = [];
  $: {
    if (process?.cpuPercent != null) {
      cpuPercents = [...cpuPercents.slice(0, 99), process.cpuPercent];
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
          <td>{component.running ? 'Running' : 'Stopped'}</td>
          <td />
        </tr>
        {#if process}
          <tr>
            <td class="label">CPU</td>
            <td>
              <PercentLabel value={process.cpuPercent} />
            </td>
          </tr>
          <tr>
            <td class="label">Resident Memory</td>
            <td>
              <BytesLabel value={process.residentBytes} />
            </td>
            <td />
          </tr>
          <tr>
            <td class="label">Started at</td>
            <td>
              <InstantLabel value={process.started} />
            </td>
            <td>
              <Sparkline entries={cpuPercents} interactive />
            </td>
          </tr>
          <tr>
            <td class="label">Local Ports</td>
            <td>{process.ports?.join(', ') ?? 'None'}</td>
            <td />
          </tr>
        {/if}
      </CheckeredTable>
      <br />
      <h3>Environment</h3>
      {#if environment}
        <EnvironmentTable variables={environment.variables} />
      {/if}
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
