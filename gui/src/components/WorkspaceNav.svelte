<script lang="ts">
  import NavbarRoute from './nav/NavbarRoute.svelte';
  import HomeSVG from './mono/HomeSVG.svelte';
  import LayersSVG from './mono/LayersSVG.svelte';
  import StorageSVG from './mono/StorageSVG.svelte';
  import VariableSVG from './mono/VariableSVG.svelte';
  import NetworkingSVG from './mono/NetworkingSVG.svelte';

  export let workspaceId: string;
  export let active: string;
</script>

<aside>
  <NavbarRoute title="Dashboard" href={`#/workspaces/${workspaceId}`} {active}>
    <HomeSVG />
  </NavbarRoute>

  <NavbarRoute
    title="Components"
    href={`#/workspaces/${workspaceId}/components`}
    {active}
  >
    <LayersSVG />
  </NavbarRoute>

  <NavbarRoute
    title="Variables"
    href={`#/workspaces/${workspaceId}/variables`}
    {active}
  >
    <VariableSVG />
  </NavbarRoute>

  <NavbarRoute
    title="Storage"
    href={`#/workspaces/${workspaceId}/storage`}
    {active}
  >
    <StorageSVG />
  </NavbarRoute>

  <NavbarRoute
    title="Networking"
    href={`#/workspaces/${workspaceId}/networking`}
    {active}
  >
    <NetworkingSVG />
  </NavbarRoute>
</aside>
