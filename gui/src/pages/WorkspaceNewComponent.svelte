<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import LayersSVG from '../components/mono/LayersSVG.svelte';
  import DockerSVG from '../components/mono/DockerSVG.svelte';
  import * as router from 'svelte-spa-router';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <Panel
    title="New component"
    backRoute={workspaceRoute}
    --panel-padding="2rem"
    --panel-overflow-y="scroll"
  >
    <div class="center-form">
      <section>
        <!-- Generic components, no heading. -->

        <button
          on:click={() => {
            router.push(
              `/workspaces/${encodeURIComponent(workspaceId)}/new-process`,
            );
          }}
        >
          <LayersSVG />
          <b>Process</b>
        </button>

        <!-- Timer, External Link, etc. -->
      </section>

      <section>
        <h2>Docker</h2>

        <button
          on:click={() => {
            router.push(
              `/workspaces/${encodeURIComponent(workspaceId)}/new-container`,
            );
          }}
        >
          <DockerSVG />
          <b>Container</b>
        </button>

        <!-- Volume, network, etc. -->
      </section>

      <!-- Databases, Apps, cloud services, etc. -->
    </div>
  </Panel>
</Layout>

<style>
  h2 {
    font-size: 20px;
    font-weight: 500;
    margin: 0;
    margin-top: 24px;
    margin-bottom: 16px;
  }

  .center-form {
    max-width: 640px;
    margin: 0 auto;
  }

  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 16px 32px 16px 24px;
    position: relative;
    display: grid;
    grid-template-columns: max-content max-content max-content;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
  }

  button:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  button:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }
</style>
