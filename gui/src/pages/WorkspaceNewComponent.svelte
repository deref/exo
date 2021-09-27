<script lang="ts">
  import Icon from '../components/Icon.svelte';
  import Layout from '../components/Layout.svelte';
  import WorkspaceNav from '../components/WorkspaceNav.svelte';
  import CenterFormPanel from '../components/form/CenterFormPanel.svelte';
  import type { IconGlyph } from '../components/Icon.svelte';
  import * as router from 'svelte-spa-router';

  export let params = { workspace: '' };

  const workspaceId = params.workspace;
  const workspaceRoute = `/workspaces/${encodeURIComponent(workspaceId)}`;

  interface ComponentType {
    displayName: string;
    name: string;
    glyph: IconGlyph;
  }

  interface Category {
    title?: string;
    componentTypes: ComponentType[];
  }

  const categories: Category[] = [
    {
      componentTypes: [
        // Generic components...
        {
          displayName: 'Process',
          name: 'process',
          glyph: 'Layers',
        },
      ],
    },
    {
      title: 'Docker',
      componentTypes: [
        // Docker components...
        {
          displayName: 'Container',
          name: 'container',
          glyph: 'Docker',
        },
        {
          displayName: 'Volume',
          name: 'volume',
          glyph: 'Docker',
        },
        {
          displayName: 'Network',
          name: 'network',
          glyph: 'Docker',
        },
      ],
    },
    // Cloud services, databases, etc...
  ];
</script>

<Layout>
  <WorkspaceNav {workspaceId} active="Dashboard" slot="navbar" />
  <CenterFormPanel title="New component" backRoute={workspaceRoute}>
    {#each categories as category}
      <section>
        {#if category.title}
          <h2>{category.title}</h2>
        {/if}
        {#each category.componentTypes as componentType}
          <button
            on:click={() => {
              router.push(
                `${workspaceRoute}/new-${encodeURIComponent(
                  componentType.name,
                )}`,
              );
            }}
          >
            <Icon glyph={componentType.glyph} />
            <b>{componentType.displayName}</b>
          </button>
        {/each}
      </section>
    {/each}
  </CenterFormPanel>
</Layout>

<style>
  section {
    margin-bottom: 24px;
  }

  button {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 16px 32px 16px 24px;
    position: relative;
    display: grid;
    width: 100%;
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
