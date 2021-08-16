<script lang="ts">
  import FormattedLogMessage from './FormattedLogMessage.svelte';
  import type { Log } from './logTypes';

  export let logs: Log[];
</script>

<table>
  {#each logs as log (log.id)}
    <tr style={log.style}>
      <td
        class="time"
        on:click={() => {
          window.alert(`Full timestamp: ${log.time.full}`);
        }}
      >
        {log.time.short}
      </td>
      <td class="name">{log.name}</td>
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
    color: #999999;
    cursor: zoom-in;
  }

  tr:hover .time {
    background: #eeeeee;
    color: #555555;
  }
</style>
