<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import BytesLabel from '../components/BytesLabel.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import EnvironmentTable from '../components/EnvironmentTable.svelte';
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';
  import sparkline from '@fnando/sparkline';
  import { api } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';
  import {
    fetchProcesses,
    refreshAllProcesses,
    processes,
  } from '../lib/process/store';
  import type { RequestLifecycle } from '../lib/api';
  import type { ProcessDescription } from '../lib/process/types';

  export let params = { workspace: '', process: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const processId = params.process;

  let processList: RequestLifecycle<ProcessDescription[]> = {
    stage: 'pending',
  };
  const unsubscribeProcesses = processes.subscribe((processes) => {
    processList = processes;
  });

  let refreshInterval: any;
  let process: ProcessDescription | null = null;

  let sparklineSvg: SVGSVGElement;

  const cpuPercentages: number[] = [];

  onMount(() => {
    fetchProcesses(workspace);

    // TODO: Server-sent events or websockets!
    refreshInterval = setInterval(() => {
      refreshAllProcesses(workspace);
      if (processList.stage === 'success') {
        process = processList.data.filter((p) => p.id === processId)[0];
      }
      if (process && process.running) {
        cpuPercentages.push(process.cpuPercent);
        if (cpuPercentages.length > 100) {
          cpuPercentages.shift();
        }
        if (
          cpuPercentages.some((p) => p !== 0) &&
          sparklineSvg &&
          sparklineSvg instanceof SVGSVGElement
        ) {
          sparkline(sparklineSvg, cpuPercentages, { interactive: true });
        }
      }
    }, 1000);
  });

  onDestroy(() => {
    clearInterval(refreshInterval);
    unsubscribeProcesses();
  });
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  {#if process}
    <Panel title={process.name} backRoute={workspaceRoute}>
      {#if process.running}
        <CheckeredTableWrapper>
          <table>
            <tbody>
              <tr>
                <td class="label">Status</td>
                <td>{process.running ? 'Running' : 'Stopped'}</td>
                <td />
              </tr>
              <tr>
                <td class="label">CPU</td>
                <td>{process.cpuPercent.toFixed(2)}%</td>
                <td
                  ><svg
                    bind:this={sparklineSvg}
                    class="sparkline"
                    width="100"
                    height="30"
                    stroke-width="3"
                  /></td
                >
              </tr>
              <tr>
                <td class="label">Resident Memory</td>
                <td><BytesLabel value={process.residentMemory} /></td>
                <td />
              </tr>
              <tr>
                <td class="label">Started at</td>
                <td
                  ><span title={new Date(process.createTime).toISOString()}
                    >{new Date(process.createTime).toLocaleTimeString()}</span
                  ></td
                >
                <td
                  ><svg
                    class="sparkline"
                    width="100"
                    height="30"
                    stroke-width="3"
                  /></td
                >
              </tr>
              <tr>
                <td class="label">Local Ports</td>
                <td>{process.ports?.join(', ') ?? 'None'}</td>
                <td />
              </tr>
              <tr>
                <td class="label">Children</td>
                <td>{process.childrenExecutables?.join(', ') ?? 'None'}</td>
                <td />
              </tr>
            </tbody>
          </table>
        </CheckeredTableWrapper>
        <br />
        <h3>Environment</h3>
        <EnvironmentTable variables={process.envVars} />
      {:else}
        <span>Process is not running</span>
      {/if}
    </Panel>
  {:else}
    <Panel title="Loading..." backRoute={workspaceRoute} />
  {/if}
</Layout>

<style>
  .label {
    font-size: 0.8em;
    font-weight: 450;
    color: var(--grey-5-color);
  }
  .sparkline {
    margin: -6px -15px;
    stroke: var(--sparkline-stroke);
    fill: var(--sparkline-fill);
  }
</style>
