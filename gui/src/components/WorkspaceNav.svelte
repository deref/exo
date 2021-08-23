<script lang="ts">
  import * as router from 'svelte-spa-router';
  import NavbarButton from './nav/NavbarButton.svelte';
  import HomeSVG from './mono/HomeSVG.svelte';
  import LayersSVG from './mono/LayersSVG.svelte';
  import StorageSVG from './mono/StorageSVG.svelte';
  import NetworkingSVG from './mono/NetworkingSVG.svelte';

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
