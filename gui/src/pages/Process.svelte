<script lang="ts">
  import * as router from 'svelte-spa-router';
  import Layout from '../components/Layout.svelte';
  import { api } from '../lib/api';
  import { onDestroy, onMount } from 'svelte';
  import { fetchProcesses, refreshAllProcesses, processes } from '../lib/process/store';
  import type { RemoteData } from '../lib/api';
  import type { ComponentDetails } from '../lib/process/types';
  import sparkline from '@fnando/sparkline';
  export let params = { workspace: '', process: '' };

  const workspaceId = params.workspace;
  const workspace = api.workspace(workspaceId);

  const processId = params.process;

  let processList: RemoteData<ComponentDetails[]> = { stage: 'pending' };
  const unsubscribeProcesses = processes.subscribe((processes) => {
    processList = processes;
  });

  let refreshInterval: any;
  let process: ComponentDetails | null = null

  const cpuPercentages: number[] = []

  onMount(() => {
    fetchProcesses(workspace);

    // TODO: Server-sent events or websockets!
    refreshInterval = setInterval(() => {
      refreshAllProcesses(workspace);
      if (processList.stage === "success") {
        process = processList.data.filter(p => p.id === processId)[0];
      }
      if(process && process.status.running) {
        cpuPercentages.push(process.status.CPUPercent)
        if (cpuPercentages.length > 100) {
          cpuPercentages.shift()
        }
        sparkline(
          document.querySelector(".sparkline")!,
          cpuPercentages,
          {interactive: true}
        );
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
          <button class="back-button" on:click={() => void router.push(`#/workspaces/${encodeURIComponent(workspaceId)}/`)}>
           🠔 Back
          </button>
          <h1>{process.name}</h1>
        </div>
        <h3>Status</h3>
        <table>
          <tr>
            <td>Status</td>
            <td>{process.status.running ? "Running" : "Stopped"}</td>
          </tr>
          <tr>
            <td>CPU</td>
            <td>{process.status.CPUPercent.toFixed(2)}%</td>
            <td><svg class="sparkline" width="100" height="30" stroke-width="3"></svg></td>
          </tr>
        </table>
        <h3>Environment</h3>
        <table>
          {#each Object.entries(process.status.envVars ?? {}) as [name, val] (name)}
            <tr>
              <td>{name}</td>
              <td>{val}</td>
            </tr>
          {/each}
        </table>
      </div>
    {:else}
      Loading...
    {/if}
  </section>
</Layout>

<style>
  section {
    height: 100%;
    display: grid;
    grid-template-columns: 360px 1fr;
    gap: 30px;
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

  td {
    padding-right: 1em;
  }

  /* line with highlight area */
  .sparkline {
    stroke: red;
    fill: rgba(255, 0, 0, .3);
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
    fill: rgba(255, 0, 0, .3);
  }

  .sparkline--line {
    stroke: red;
  }
</style>
