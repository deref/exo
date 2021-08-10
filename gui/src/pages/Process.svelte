<script lang="ts">
  import sparkline from '@fnando/sparkline';
  import Layout from '../components/Layout.svelte';
  import { api } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';
  import {
    fetchProcesses,
    refreshAllProcesses,
    processes,
  } from '../lib/process/store';
  import type { RemoteData } from '../lib/api';
  import BytesLabel from '../components/BytesLabel.svelte';
  import WithLeftWorkspaceNav from '../components/WithLeftWorkspaceNav.svelte';
  import CheckeredTableWrapper from '../components/CheckeredTableWrapper.svelte';
  import type { ProcessDescription } from 'src/lib/process/types';
  export let params = { workspace: '', process: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const processId = params.process;

  let processList: RemoteData<ProcessDescription[]> = { stage: 'pending' };
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

<Layout showBackButton backButtonRoute={workspaceRoute}>
  <WithLeftWorkspaceNav {workspaceId} active="Dashboard">
    <section>
      {#if process}
        <div>
          <div id="heading">
            <h1>{process.name}</h1>
          </div>
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
                        >{new Date(
                          process.createTime,
                        ).toLocaleTimeString()}</span
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
            <CheckeredTableWrapper>
              <tbody>
                <table>
                  {#each Object.entries(process.envVars ?? {}) as [name, val] (name)}
                    <tr>
                      <td class="label">{name}</td>
                      <td><code><pre>{val}</pre></code></td>
                    </tr>
                  {/each}
                </table>
              </tbody>
            </CheckeredTableWrapper>
            <br />
          {:else}
            <p>Process is not running</p>
          {/if}
        </div>
      {:else}
        Loading...
      {/if}
    </section>
  </WithLeftWorkspaceNav>
</Layout>

<style>
  section {
    height: 100%;
    margin: 0 30px;
    padding-bottom: 16px;
  }

  #heading {
    display: flex;
    align-items: center;
  }

  .sparkline {
    stroke: red;
    fill: none;
    margin: -6px -15px;
  }

  code {
    width: 100%;
    max-width: 600px;
    display: inline-block;
    overflow-x: auto;
    font-size: 1.05em;
    padding: 8px;
    margin: -10px;
  }

  .label {
    font-size: 0.8em;
    font-weight: 450;
    color: #555555;
  }

  /* line with highlight area */
  .sparkline {
    stroke: red;
    fill: rgba(255, 0, 0, 0.3);
  }
</style>
