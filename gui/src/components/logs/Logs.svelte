<script lang="ts">
  import { afterUpdate, beforeUpdate } from 'svelte';

  import FormattedLogMessage from './FormattedLogMessage.svelte';
  import type { LogEvent } from './types';

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

    const threshold = 150;
    const currentPosition = logViewport.scrollTop + logViewport.offsetHeight;
    const height = logViewport.scrollHeight;
    wasScrolledCloseToBottom = currentPosition > height - threshold;
  });

  afterUpdate(async () => {
    if (wasScrolledCloseToBottom && logViewport) {
      logViewport.scrollTop = logViewport.scrollHeight;
    }
  });
</script>

<div class="logs-container" bind:this={logViewport}>
  <table>
    {#each events as event (event.id)}
      <tr style={event.style}>
        <td
          class="time"
          on:click={() => {
            window.alert(`Full timestamp: ${event.time.full}`);
          }}
        >
          {event.time.short}
        </td>
        <td class="name">{event.name}</td>
        <td>
          <FormattedLogMessage message={event.message} />
        </td>
      </tr>
    {/each}
  </table>
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

  table,
  tr,
  td {
    border: none;
    border-collapse: collapse;
  }

  td {
    padding: 0 0.3em;
    vertical-align: text-top;
    color: var(--grey-3-color);
    white-space: pre-wrap;
  }

  tr:hover td {
    background: var(--grey-e-color);
    color: var(--grey-1-color);
  }

  .name {
    text-align: right;
    background: var(--log-bg-color);
    color: var(--log-color);
  }

  tr:hover .name {
    background: var(--log-bg-hover-color);
    color: var(--log-hover-color);
  }

  @media (prefers-color-scheme: dark) {
    .name {
      background: var(--dark-log-bg-color);
      color: var(--dark-log-color);
    }

    tr:hover .name {
      background: var(--dark-log-bg-hover-color);
      color: var(--dark-log-hover-color);
    }
  }

  .time {
    color: var(--grey-9-color);
    cursor: zoom-in;
  }

  tr:hover .time {
    color: var(--grey-5-color);
  }
</style>
