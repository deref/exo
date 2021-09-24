<script lang="ts">
  import FormattedLogMessage from './FormattedLogMessage.svelte';
  import type { LogEvent } from './types';

  export let event: LogEvent;

  // SEE NOTE [LOG_COMPONENTS]
  const trimLabel = (name: string) => name.replace(/(:err$)|(:out$)/g, '');
</script>

<tr style={event.style}>
  <td
    class="time"
    on:click={() => {
      window.alert(`Full timestamp: ${event.time.full}`);
    }}
  >
    {event.time.short}
  </td>
  <td class="name" title={event.name}>{trimLabel(event.name)}</td>
  <td>
    <FormattedLogMessage message={event.message} />
  </td>
</tr>

<style>
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

  .time {
    color: var(--grey-9-color);
    cursor: zoom-in;
  }

  tr:hover .time {
    color: var(--grey-5-color);
  }
</style>
