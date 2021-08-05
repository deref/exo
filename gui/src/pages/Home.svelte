<script lang="ts">
  import Layout from '../components/Layout.svelte';
  import Code from '../components/Code.svelte';
  import * as router from 'svelte-spa-router';
  import { api } from '../lib/api';

  const workspaces = api.kernel.describeWorkspaces();
</script>

<Layout>
  <section>
    Use <Code>exo gui</Code> in your terminal to launch into the current directory's
    workspace.
    <h2>Workspaces</h2>
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
              <span>{workspace.root}</span>
            </li>
          {/each}
        </ul>
      {:catch error}
        <p style="color: red">{error.message}</p>
      {/await}
    </div>
  </section>
</Layout>

<style>
  section {
    background: #ffffff;
    padding: 30px;
    min-height: 100%;
  }

  h2 {
    margin-top: 36px;
    margin-bottom: 24px;
  }

  ul {
    list-style: none;
    margin: 0;
    padding: 0;
    display: grid;
    gap: 16px;
    grid-template-columns: repeat(
      auto-fill,
      minmax(320px, 1fr)
    ); /* Auto-responsive */
  }

  li {
    background: linear-gradient(#fff, #f5f5f5);
    border: none;
    border-radius: 6px;
    padding: 16px 24px;
    box-shadow: 0 0.33px 0 1px hsla(0, 0%, 100%, 0.15),
      0 4px 8px -3px rgba(0, 0, 0, 0.15), 0 0.4px 0 0.8px rgba(0, 0, 0, 0.25);
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
    background: linear-gradient(#fafafa, #e7e7e7);
    box-shadow: 0 0.33px 0 1px hsla(0, 0%, 100%, 0.15),
      0 6px 8px -4px rgba(0, 0, 0, 0.2), 0 0.4px 0 0.8px rgba(0, 0, 0, 0.35);
  }

  li:active {
    background: linear-gradient(#f7f7f7, #e0e0e0);
    box-shadow: 0 0.33px 0 1px hsla(0, 0%, 100%, 0.15),
      0 4px 6px -3px rgba(0, 0, 0, 0.15), 0 0.4px 0 0.8px rgba(0, 0, 0, 0.45);
  }
</style>
