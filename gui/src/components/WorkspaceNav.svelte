<script lang="ts">
  import Icon from './Icon.svelte';
  import NavbarRoute from './nav/NavbarRoute.svelte';

  export let workspaceId: string;
  export let active: string;
</script>

<aside>
  <NavbarRoute title="Dashboard" href={`#/workspaces/${workspaceId}`} {active}>
    <Icon glyph="Home" />
  </NavbarRoute>

  <NavbarRoute
    title="Components"
    href={`#/workspaces/${workspaceId}/components`}
    {active}
  >
    <Icon glyph="Layers" />
  </NavbarRoute>

  <NavbarRoute
    title="Variables"
    href={`#/workspaces/${workspaceId}/variables`}
    {active}
  >
    <Icon glyph="Variable" />
  </NavbarRoute>

  <NavbarRoute
    title="Storage"
    href={`#/workspaces/${workspaceId}/storage`}
    {active}
  >
    <Icon glyph="Storage" />
  </NavbarRoute>

  <NavbarRoute
    title="Networking"
    href={`#/workspaces/${workspaceId}/networking`}
    {active}
  >
    <Icon glyph="Networking" />
  </NavbarRoute>
</aside>
