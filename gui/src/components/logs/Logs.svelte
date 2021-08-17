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

  .time {
    color: var(--grey-9-color);
    cursor: zoom-in;
  }

  tr:hover .time {
    color: var(--grey-5-color);
  }
</style>
