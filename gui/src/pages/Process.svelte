<script lang="ts">
  import * as router from 'svelte-spa-router';
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
  import type { ComponentDetails } from '../lib/process/types';
  import BytesLabel from '../components/BytesLabel.svelte';
  export let params = { workspace: '', process: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);

  const processId = params.process;

  let processList: RemoteData<ComponentDetails[]> = { stage: 'pending' };
  const unsubscribeProcesses = processes.subscribe((processes) => {
    processList = processes;
  });

  let refreshInterval: any;
  let process: ComponentDetails | null = null;

  const cpuPercentages: number[] = [];

  onMount(() => {
    fetchProcesses(workspace);

    // TODO: Server-sent events or websockets!
    refreshInterval = setInterval(() => {
      refreshAllProcesses(workspace);
      if (processList.stage === 'success') {
        process = processList.data.filter((p) => p.id === processId)[0];
      }
      if (process && process.status.running) {
        cpuPercentages.push(process.status.cpuPercent);
        if (cpuPercentages.length > 100) {
          cpuPercentages.shift();
        }
        const sparklineSvg = document.querySelector('.sparkline');
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
  <section>
    {#if process}
      <div>
        <div id="heading">
          <button
            class="back-button"
            on:click={() =>
              void router.push(
                `#/workspaces/${encodeURIComponent(workspaceId)}/`,
              )}
          >
            🠔 Back
          </button>
          <h1>{process.name}</h1>
        </div>
        <h3>Status</h3>
        {#if process.status.running}
          <table>
            <tr>
              <td>Status</td>
              <td>{process.status.running ? 'Running' : 'Stopped'}</td>
            </tr>
            <tr>
              <td>CPU</td>
              <td>{process.status.cpuPercent.toFixed(2)}%</td>
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
              <td>Resident Memory</td>
              <td><BytesLabel value={process.status.residentMemory} /></td>
            </tr>
            <tr>
              <td>Started at</td>
              <td
                ><span title={new Date(process.status.createTime).toISOString()}
                  >{new Date(
                    process.status.createTime,
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
              <td>Local Ports</td>
              <td>{process.status.ports?.join(', ') ?? 'None'}</td>
            </tr>
            <tr>
              <td>Children</td>
              <td>{process.status.childrenExecutables?.join(', ') ?? 'None'}</td
              >
            </tr>
          </table>
          <h3>Environment</h3>
          <table>
            {#each Object.entries(process.status.envVars ?? {}) as [name, val] (name)}
              <tr>
                <td>{name}</td>
                <td><code><pre>{val}</pre></code></td>
              </tr>
            {/each}
          </table>
        {:else}
          <p>Process is not running</p>
        {/if}
      </div>
    {:else}
      Loading...
    {/if}
  </section>
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

  .back-button {
    padding: 0.5em;
    padding-bottom: 0.3em;
    margin-right: 1em;
  }

  .sparkline {
    stroke: red;
    fill: none;
  }

  code {
    width: 100%;
    max-width: 500px;
    display: inline-block;
    overflow-x: auto;
    padding: 0.6em;
    border-radius: 0.5em;
    background-color: rgba(0, 0, 0, 0.05);
  }

  td {
    padding-right: 2em;
  }

  /* line with highlight area */
  .sparkline {
    stroke: red;
    fill: rgba(255, 0, 0, 0.3);
  }

  /* change the spot color */
  .sparkline--spot {
    stroke: blue;
    fill: blue;
  }

  /* change the cursor color */
  .sparkline--cursor {
    stroke: orange;
  }

  /* style fill area and line colors using specific class name */
  .sparkline--fill {
    fill: rgba(255, 0, 0, 0.3);
  }

  .sparkline--line {
    stroke: red;
  }
</style>
