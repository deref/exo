<script lang="ts">
  import * as router from 'svelte-spa-router';

  import Home from './mono/home.svelte';
  import Layers from './mono/layers.svelte';
  import Storage from './mono/database.svelte';
  import Networking from './mono/networking.svelte';

  export let workspaceId: string;
  export let active: string;
</script>

<div>
  <aside>
    <button
      on:click={() => {
        router.push(`#/workspaces/${workspaceId}`);
      }}
      class={active === 'Dashboard' ? 'active' : ''}
      title="Dashboard"><Home /></button
    >
    <button
      on:click={() => {
        router.push(`#/workspaces/${workspaceId}/components`);
      }}
      class={active === 'Components' ? 'active' : ''}
      title="Components"><Layers /></button
    >
    <button
      on:click={() => {
        router.push(`#/workspaces/${workspaceId}/storage`);
      }}
      class={active === 'Storage' ? 'active' : ''}
      title="Storage"><Storage /></button
    >
    <button
      on:click={() => {
        router.push(`#/workspaces/${workspaceId}/networking`);
      }}
      class={active === 'Networking' ? 'active' : ''}
      title="Networking"><Networking /></button
    >
  </aside>
  <section>
    <slot />
  </section>
</div>

<style>
  div {
    position: relative;
    overflow: hidden;
    display: grid;
    grid-template-columns: max-content auto;
    gap: 1px;
    height: 100%;
  }

  div > section {
    position: relative;
    overflow-y: auto;
    overflow-x: hidden;
    background: #ffffff;
  }

  aside {
    width: 48px;
    background: #e9e9e9;
  }

  button {
    position: relative;
    border: none;
    background: #e9e9e9;
    color: #999999;
    height: 48px;
    width: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  button > :global(svg),
  button > :global(svg *) {
    fill: currentColor;
  }

  button.active {
    background: #d5d5d5;
    color: #333333;
  }

  button:hover {
    background: #dddddd;
    color: #666666;
  }

  button.active:hover {
    background: #d0d0d0;
    color: #222222;
  }
</style>
