<script lang="ts">
  import * as router from 'svelte-spa-router';

  import Home from './mono/home.svelte';
  import Layers from './mono/layers.svelte';
  import Storage from './mono/database.svelte';
  import Networking from './mono/networking.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';

  export let workspaceId: string;
  export let active: string;
</script>

<aside>
  <NavbarButton
    title="Dashboard"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}`);
    }}
    {active}
  >
    <Home />
  </NavbarButton>

  <NavbarButton
    title="Components"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}/components`);
    }}
    {active}
  >
    <Layers />
  </NavbarButton>

  <NavbarButton
    title="Storage"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}/storage`);
    }}
    {active}
  >
    <Storage />
  </NavbarButton>

  <NavbarButton
    title="Networking"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}/networking`);
    }}
    {active}
  >
    <Networking />
  </NavbarButton>
</aside>

<style>
  aside {
    flex-grow: 1;
  }
</style>
