<script lang="ts">
  import { afterUpdate, beforeUpdate } from 'svelte';

  import LogRow from './LogRow.svelte';
  import Tooltip from '../Tooltip.svelte';
  import type { GetComponentNameFunc } from './LogRow.svelte';
  import type { LogEvent } from '../../lib/logs/types';

  export let getComponentName: GetComponentNameFunc | undefined;
  export let events: LogEvent[] = [];

  // Automatically scroll on new logs if the user is already scrolled close to the bottom of the content.
  let logViewport: HTMLElement;
  let wasScrolledCloseToBottom = true;
  // Record whether the user was scrolled close to the bottom before new entries arrived.
  // If so, scroll them to the new bottom after the update.
  beforeUpdate(() => {
    if (!logViewport) {
      return;
    }

    const threshold = 20; // Approximate minimum height of line log line.
    const currentPosition = logViewport.scrollTop + logViewport.offsetHeight;
    const height = logViewport.scrollHeight;
    wasScrolledCloseToBottom = currentPosition > height - threshold;
  });

  const scrollToBottom = () => {
    logViewport.scrollTop = logViewport.scrollHeight;
  };

  afterUpdate(async () => {
    if (wasScrolledCloseToBottom && logViewport) {
      scrollToBottom();
    }
  });
</script>

<div class="logs-container" bind:this={logViewport}>
  <table>
    {#each events as event (event.id)}
      <LogRow {getComponentName} {event} />
    {/each}
  </table>
  <button
    class="latest"
    class:showLatest={!wasScrolledCloseToBottom}
    on:click={(e) => {
      scrollToBottom();
    }}
  >
    <Tooltip>↓ Latest events</Tooltip>
  </button>
</div>

<style>
  .logs-container {
    width: 100%;
    height: 100%;
    overflow-y: scroll;
    overflow-x: hidden;
  }

  table {
    background: var(--primary-bg-color);
    font-family: var(--font-mono);
    font-variant-ligatures: var(--preferred-ligatures-logs);
    font-weight: 450;
    font-size: 15px;
  }

  table {
    border: none;
    border-collapse: collapse;
  }

  .latest {
    display: none;
    position: absolute;
    right: 16px;
    bottom: 12px;
    cursor: pointer;
    margin: 0;
    padding: 0;
    border: none;
    background: none;
  }

  .showLatest {
    display: block;
  }
</style>
