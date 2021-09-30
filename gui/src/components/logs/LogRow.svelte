<script lang="ts" context="module">
  export type GetComponentNameFunc = (id: string) => string | null;
</script>

<script lang="ts">
  import FormattedLogMessage from './FormattedLogMessage.svelte';
  import { shortTime } from '../../lib/time';
  import type { LogEvent } from '../../lib/logs/types';
  import { logStyleFromHash } from '../../lib/color';

  export let getComponentName: GetComponentNameFunc = (id) => null;
  export let event: LogEvent;

  const componentName = getComponentName(event.stream) ?? event.stream;
</script>

<tr style={logStyleFromHash(componentName)}>
  <td
    class="time"
    on:click={() => {
      window.alert(`Full timestamp: ${event.timestamp}`);
    }}
  >
    {shortTime(event.timestamp)}
  </td>
  <td class="name" title={event.stream}>
    {componentName}
  </td>
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
