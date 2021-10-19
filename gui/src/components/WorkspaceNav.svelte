<script lang="ts">
  import Icon from './Icon.svelte';
  import NavbarButton from './nav/NavbarButton.svelte';
  import type { IconGlyph } from './Icon.svelte';

  export let workspaceId: string;
  export let active: string;

  interface NavRoute {
    title: string;
    href: string;
    glyph: IconGlyph;
  }

  const routes: NavRoute[] = [
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
    <NavbarButton title={route.title} href={route.href} {active}>
      <Icon glyph={route.glyph} />
    </NavbarButton>
  {/each}
</aside>
