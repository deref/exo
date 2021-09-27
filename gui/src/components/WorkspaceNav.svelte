<script lang="ts">
  import Icon from './Icon.svelte';
  import NavbarRoute from './nav/NavbarRoute.svelte';

  export let workspaceId: string;
  export let active: string;

  const routes = [
    {
      title: 'Dashboard',
      href: `#/workspaces/${workspaceId}`,
      glyph: 'Home',
    },
    {
      title: 'Components',
      href: `#/workspaces/${workspaceId}/components`,
      glyph: 'Layers',
    },
    {
      title: 'Variables',
      href: `#/workspaces/${workspaceId}/variables`,
      glyph: 'Variable',
    },
    {
      title: 'Storage',
      href: `#/workspaces/${workspaceId}/storage`,
      glyph: 'Storage',
    },
    {
      title: 'Networking',
      href: `#/workspaces/${workspaceId}/networking`,
      glyph: 'Networking',
    },
  ];
</script>

<aside>
  {#each routes as route}
    <NavbarRoute title={route.title} href={route.href} {active}>
      <Icon glyph={route.glyph} />
    </NavbarRoute>
  {/each}
</aside>
