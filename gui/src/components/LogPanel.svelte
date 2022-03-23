<script lang="ts" context="module">
  import type { Readable } from 'svelte/store';
  import type { Event } from './logs/LogScroller.svelte';

  export type Stream = {
    events: Event[];
    filterString: string;
  };

  export type { Event };

  export type StreamStore = Readable<Stream> & {
    clearEvents(): void;
    setFilterString(value: string): void;
  };
</script>

<script lang="ts">
  import Panel from './Panel.svelte';
  import LogScroller from './logs/LogScroller.svelte';
  import LogFilterBar from './logs/LogFilterBar.svelte';
  import { derived } from 'svelte/store';

  export let stream: StreamStore;

  const filterString = {
    ...derived(stream, ({ filterString }) => filterString),
    set: (value: string) => {
      stream.setFilterString(value);
    },
  };
</script>

<Panel title="Logs" --panel-padding="0" --panel-overflow-y="hidden">
  {#if $stream.events.length}
    <LogScroller events={$stream.events} />
  {:else}
    <div class="empty">
      {#if $filterString}
        No events match: "<code>{$filterString}</code>"
      {:else}
        No events
      {/if}
    </div>
  {/if}
  <LogFilterBar
    slot="bottom"
    bind:filterString={$filterString}
    clearEvents={stream.clearEvents}
  />
</Panel>

<style>
  .empty {
    font-style: italic;
    padding: 8px;
    font-size: 15px;
  }

  .empty code {
    font-style: normal;
    font-weight: 500;
    white-space: pre-wrap;
  }
</style>
