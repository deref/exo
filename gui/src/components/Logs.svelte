<script lang="ts">
  import FormattedLogMessage from './logs/FormattedLogMessage.svelte';
  import type { Logs } from './logs/logTypes';

  export let logs: Logs;
</script>

<table>
  {#each logs as log (log.id)}
    <tr class="log-entry" style={log.style}>
      <td class="timestamp">
        <span class="short-time">{log.time.short}</span>
        <span class="full-timestamp">{log.time.full}</span>
      </td>
      <td>{log.name}</td>
      <td>
        <FormattedLogMessage message={log.message} />
      </td>
    </tr>
  {/each}
</table>

<style>
  table {
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
    color: #333333;
    white-space: pre-wrap;
  }

  tr:hover td {
    background: #f3f3f3;
    color: #111111;
  }

  td:nth-child(1) {
    color: #999999;
  }

  tr:hover td:nth-child(1) {
    background: #eeeeee;
    color: #555555;
  }

  td:nth-child(2) {
    text-align: right;
    background: var(--log-bg-color);
    color: var(--log-color);
  }

  tr:hover td:nth-child(2) {
    background: var(--log-bg-hover-color);
    color: var(--log-hover-color);
  }

  .timestamp {
    position: relative;
  }

  .full-timestamp {
    display: none;
  }

  .timestamp:hover {
    overflow: visible;
  }

  .timestamp:hover .full-timestamp {
    display: block;
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    width: auto;
    white-space: nowrap;
    background: #333333;
    color: #f3f3f3;
    padding: 0 0.3em;
  }
</style>
