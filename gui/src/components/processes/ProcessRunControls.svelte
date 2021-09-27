<script lang="ts">
  import IconButton from '../IconButton.svelte';
  import Spinner from '../Spinner.svelte';
  import PlaySVG from '../mono/PlaySVG.svelte';
  import PauseSVG from '../mono/PauseSVG.svelte';

  export let setProcRun: (id: string, run: boolean) => void;
  export let statusPending: Set<string>;
  export let id: string;
  export let running: boolean;
</script>

<div class="run-controls">
  {#if statusPending.has(id)}
    <Spinner />
  {:else if running}
    <div class="running unhover-only" />
    <div class="control hover-only">
      <IconButton tooltip="Stop process" on:click={() => setProcRun(id, false)}>
        <PauseSVG />
      </IconButton>
    </div>
  {:else}
    <div class="stopped unhover-only" />
    <div class="control hover-only">
      <IconButton tooltip="Run process" on:click={() => setProcRun(id, true)}>
        <PlaySVG />
      </IconButton>
    </div>
  {/if}
</div>

<style>
  :global(div:not(:hover):not(:focus-within) > div)
    > .run-controls
    .hover-only {
    display: none;
  }

  :global(div:hover > div) > .run-controls .unhover-only,
  :global(div:focus-within > div) > .run-controls .unhover-only {
    display: none;
  }

  .run-controls {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
  }

  .stopped {
    width: 16px;
    height: 16px;
    border-radius: 8px;
    border: 1px solid var(--grey-a-color);
  }

  .running {
    width: 16px;
    height: 16px;
    border-radius: 8px;
    background: var(--online-green-color);
  }

  .control {
    position: absolute;
    z-index: 3;
  }
</style>
