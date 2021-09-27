<script lang="ts" context="module">
  // Note, these are internal named SVG icons. Dynamic SVGs should be handled separately.

  // This file should be generated in a similar fashion to the way the theme file is generated.

  export type IconGlyph =
    | 'Add'
    | 'Back'
    | 'Check'
    | 'Delete'
    | 'Details'
    | 'Docker'
    | 'Ellipsis'
    | 'Feedback'
    | 'Home'
    | 'Layers'
    | 'LeftUp'
    | 'Logs'
    | 'Networking'
    | 'Pause'
    | 'Play'
    | 'Preferences'
    | 'Reset'
    | 'Storage'
    | 'Variable'
    | 'Warning';
</script>

<script lang="ts">
  import AddSVG from './mono/AddSVG.svelte';
  import BackSVG from './mono/BackSVG.svelte';
  import CheckSVG from './mono/CheckSVG.svelte';
  import DeleteSVG from './mono/DeleteSVG.svelte';
  import DetailsSVG from './mono/DetailsSVG.svelte';
  import DockerSVG from './mono/DockerSVG.svelte';
  import EllipsisSVG from './mono/EllipsisSVG.svelte';
  import FeedbackSVG from './mono/FeedbackSVG.svelte';
  import HomeSVG from './mono/HomeSVG.svelte';
  import LayersSVG from './mono/LayersSVG.svelte';
  import LeftUpSVG from './mono/LeftUpSVG.svelte';
  import LogsSVG from './mono/LogsSVG.svelte';
  import NetworkingSVG from './mono/NetworkingSVG.svelte';
  import PauseSVG from './mono/PauseSVG.svelte';
  import PlaySVG from './mono/PlaySVG.svelte';
  import PreferencesSVG from './mono/PreferencesSVG.svelte';
  import ResetSVG from './mono/ResetSVG.svelte';
  import StorageSVG from './mono/StorageSVG.svelte';
  import VariableSVG from './mono/VariableSVG.svelte';
  import WarningSVG from './mono/WarningSVG.svelte';

  export let glyph: IconGlyph;
</script>

{#if glyph === 'Add'}
  <AddSVG />
{:else if glyph === 'Back'}
  <BackSVG />
{:else if glyph === 'Check'}
  <CheckSVG />
{:else if glyph === 'Delete'}
  <DeleteSVG />
{:else if glyph === 'Details'}
  <DetailsSVG />
{:else if glyph === 'Docker'}
  <DockerSVG />
{:else if glyph === 'Ellipsis'}
  <EllipsisSVG />
{:else if glyph === 'Feedback'}
  <FeedbackSVG />
{:else if glyph === 'Home'}
  <HomeSVG />
{:else if glyph === 'Layers'}
  <LayersSVG />
{:else if glyph === 'LeftUp'}
  <LeftUpSVG />
{:else if glyph === 'Logs'}
  <LogsSVG />
{:else if glyph === 'Networking'}
  <NetworkingSVG />
{:else if glyph === 'Pause'}
  <PauseSVG />
{:else if glyph === 'Play'}
  <PlaySVG />
{:else if glyph === 'Preferences'}
  <PreferencesSVG />
{:else if glyph === 'Reset'}
  <ResetSVG />
{:else if glyph === 'Storage'}
  <StorageSVG />
{:else if glyph === 'Variable'}
  <VariableSVG />
{:else if glyph === 'Warning'}
  <WarningSVG />
{:else}
  <LayersSVG />
{/if}
