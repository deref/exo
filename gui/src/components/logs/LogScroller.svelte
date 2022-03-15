<script lang="ts" context="module">
  import type { Event } from './LogRow.svelte';

  export type { Event };
</script>

<script lang="ts">
  import { afterUpdate, beforeUpdate } from 'svelte';

  import LogRow from './LogRow.svelte';

  export let events: Event[] = [];

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

<div class="container" bind:this={logViewport}>
  <table>
    {#each events as event (event.id)}
      <LogRow {event} />
    {/each}
  </table>
  <div class="latest" class:showLatest={!wasScrolledCloseToBottom}>
    <button
      on:click={(e) => {
        scrollToBottom();
      }}
    >
      ↓ Latest events
    </button>
  </div>
</div>

<style>
  .container {
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
    margin: 10;
    padding: 0;
  }

  .latest button {
    cursor: pointer;
    border: none;
    background: none;
    display: inline-block;
    font-size: 13px;
    font-weight: 400;
    color: var(--strong-color);
    background: var(--primary-bg-color);
    box-shadow: var(--card-shadow);
    height: 22px;
    padding: 3px 6px;
    border-radius: 4px;
    white-space: nowrap;
  }

  .latest button:hover {
    background: var(--grey-e-color);
  }

  .showLatest {
    display: block;
  }
</style>
