<script lang="ts">
  import IconButton from '../IconButton.svelte';
  import RunSVG from '../mono/play.svelte';
  import StopSVG from '../mono/stop.svelte';

  export let setProcRun: (id: string, run: boolean) => void;
  export let statusPending: Set<string>;
  export let id: string;
  export let running: boolean;
</script>

<div class="run-controls">
  {#if statusPending.has(id)}
    <div class="spinner" />
  {:else if running}
    <div class="running unhover-only" />
    <div class="control hover-only">
      <IconButton tooltip="Stop process" on:click={() => setProcRun(id, false)}>
        <StopSVG />
      </IconButton>
    </div>
  {:else}
    <div class="stopped unhover-only" />
    <div class="control hover-only">
      <IconButton tooltip="Run process" on:click={() => setProcRun(id, true)}>
        <RunSVG />
      </IconButton>
    </div>
  {/if}
</div>

<style>
  :global(tr:not(:hover):not(:focus-within) > td) > .run-controls .hover-only {
    display: none;
  }

  :global(tr:hover > td) > .run-controls .unhover-only,
  :global(tr:focus-within > td) > .run-controls .unhover-only {
    display: none;
  }

  .run-controls {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    margin-right: 8px;
  }

  .stopped {
    width: 14px;
    height: 14px;
    border-radius: 2px;
    background: var(--grey-c-color);
  }

  .running {
    width: 16px;
    height: 16px;
    border-radius: 8px;
    background: var(--online-green-color);
  }

  .control {
    position: absolute;
    z-index: 4;
  }

  .spinner {
    position: absolute;
    z-index: 3;
    top: 7px;
    left: 7px;
    width: 18px;
    height: 18px;
    border-radius: 100%;
    animation: spin 1s infinite linear;
    border: 2px solid;
    border-top-color: var(--spinner-grey-t);
    border-right-color: var(--spinner-grey-r);
    border-bottom-color: var(--spinner-grey-b);
    border-left-color: var(--spinner-grey-l);
    transition: all 0.125s;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>
