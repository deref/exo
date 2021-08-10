<script lang="ts">
  import * as router from 'svelte-spa-router';
  import NavbarButton from './nav/NavbarButton.svelte';
  import HomeSVG from './mono/home.svelte';
  import LayersSVG from './mono/layers.svelte';
  import StorageSVG from './mono/database.svelte';
  import NetworkingSVG from './mono/networking.svelte';

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
    <HomeSVG />
  </NavbarButton>

  <NavbarButton
    title="Components"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}/components`);
    }}
    {active}
  >
    <LayersSVG />
  </NavbarButton>

  <NavbarButton
    title="Storage"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}/storage`);
    }}
    {active}
  >
    <StorageSVG />
  </NavbarButton>

  <NavbarButton
    title="Networking"
    on:click={() => {
      router.push(`#/workspaces/${workspaceId}/networking`);
    }}
    {active}
  >
    <NetworkingSVG />
  </NavbarButton>
</aside>

<style>
  aside {
    flex-grow: 1;
  }
</style>
