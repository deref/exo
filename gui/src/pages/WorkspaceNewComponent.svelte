<script lang="ts">
  import Panel from '../components/Panel.svelte';
  import Layout from '../components/Layout.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import LayersSVG from '../components/mono/LayersSVG.svelte';
  import * as router from 'svelte-spa-router';
  // import { api } from '../lib/api';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  // const workspace = api.workspace(workspaceId);
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  const componentTypes = [
    {
      type: 'process',
      title: 'Process',
      subtitle: null,
    },
  ];
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <Panel title="New component" backRoute={workspaceRoute}>
    {#each componentTypes as component}
      <button
        on:click={() => {
          router.push(
            `/workspaces/${encodeURIComponent(
              workspaceId,
            )}/new-${encodeURIComponent(component.type)}`,
          );
        }}
      >
        <LayersSVG />
        <b>{component.title}</b>
        {#if component.subtitle}
          <span>{component.subtitle}</span>
        {/if}
      </button>
    {/each}
  </Panel>
</Layout>

<style>
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
    margin-top: 12px;
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
