<script lang="ts">
  import Code from '../components/Code.svelte';
  import Layout from '../components/Layout.svelte';
  import Panel from '../components/Panel.svelte';
  import * as router from 'svelte-spa-router';
  import { api } from '../lib/api';

  const workspaces = api.kernel
    .describeWorkspaces()
    .then((workspaces) =>
      workspaces.sort((w1, w2) => w1.root.localeCompare(w2.root)),
    );
</script>

<Layout>
  <Panel title="Workspaces">
    <p>
      Use <Code>exo gui</Code> in your terminal to launch into the current directory's
      workspace.
    </p>
    <div>
      {#await workspaces}
        loading workspaces...
      {:then workspaces}
        <ul>
          {#each workspaces as workspace}
            <li
              class="a"
              on:click={() => {
                router.push(`/workspaces/${encodeURIComponent(workspace.id)}`);
              }}
            >
              <b>{workspace.id}</b>
              <span title={workspace.root}>{workspace.root}</span>
            </li>
          {/each}
        </ul>
      {:catch error}
        <p style="color: red">{error.message}</p>
      {/await}
    </div>
  </Panel>
</Layout>

<style>
  p {
    margin: 0;
    margin-bottom: 24px;
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 16px;
    /* Auto-responsive */
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  }

  li {
    background: var(--button-background);
    box-shadow: var(--button-shadow);
    border: none;
    border-radius: 4px;
    padding: 16px 24px;
    position: relative;
    display: grid;
    grid-template-columns: max-content 2fr;
    align-items: center;
    gap: 12px;
  }

  li > b {
    max-width: 6em;
  }

  li > * {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  li:hover {
    background: var(--button-hover-background);
    box-shadow: var(--button-hover-shadow);
  }

  li:active {
    background: var(--button-active-background);
    box-shadow: var(--button-active-shadow);
  }
</style>
