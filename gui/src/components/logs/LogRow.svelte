<script lang="ts" context="module">
  export type Event = {
    id: string;
    sourceId: string;
    sourceName: string;
    timestamp: string;
    message: string;
  };
</script>

<script lang="ts">
  import FormattedLogMessage from './FormattedLogMessage.svelte';
  import { shortTime } from '../../lib/time';
  import { logStyleFromHash } from '../../lib/color';

  import { modal } from '../../lib/modal';
  import { bind } from '../modal/Modal.svelte';
  import ModalDefaultPopup from '../modal/ModalDefaultPopup.svelte';

  export let event: Event;

  const showFullTimeModal = () => {
    modal.set(
      bind(ModalDefaultPopup, {
        title: shortTime(event.timestamp),
        message: `Full timestamp: ${event.timestamp}`,
      }),
    );
  };
</script>

<tr style={logStyleFromHash(event.sourceId)}>
  <td class="time" on:click={showFullTimeModal}>
    {shortTime(event.timestamp)}
  </td>
  <td class="name" title={event.sourceName}>
    {event.sourceName}
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
    white-space: nowrap;
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
