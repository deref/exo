<script lang="ts">
  import Panel from './Panel.svelte';
  import Logs from './logs/Logs.svelte';
  import LocalLogProvider from './logs/LocalLogProvider.svelte';
  import type { WorkspaceApi } from '../lib/api';
  import debounce from '../lib/debounce';
  import { visibleLogsStore } from '../lib/logs/visible-logs';
  import { setsIdentical } from '../lib/sets';
  import { processes } from '../lib/process/store';
  import { onDestroy } from 'svelte';

  export let workspace: WorkspaceApi;

  // Track which logs are marked visible.
  let logs: string[] = [];
  visibleLogsStore.subscribe((logsSet) => {
    // The `visibleLogsStore` is derived from `processes`, so it will update every
    // every time the processes are refreshed (which happens on an interval). To
    // minimize requests/flicket, we only update the logs when there is an actual change.
    if (!setsIdentical(logsSet, new Set(logs))) {
      logs = [...logsSet.values()];
    }
  });

  // Maintain a mapping of process id to name. This can be removed once log streams are
  // tagged with a process name.
  let processIdToName: Record<string, string> = {};
  const unsubscribeProcessStore = processes.subscribe((processData) => {
    if (processData.stage === 'success' || processData.stage === 'refetching') {
      processIdToName = processData.data.reduce(
        (acc, processDescription) => ({
          ...acc,
          [processDescription.id]: processDescription.name,
        }),
        {},
      );
    }
  });

  // Update filterStr when the input state changes (debounced).
  let filterStr = '';
  const setFilterStrDebounced = debounce((s: string) => {
    if (filterStr !== s) {
      filterStr = s;
    }
  }, 250);
  let filterInput = '';
  $: {
    if (filterInput !== null) {
      const text = filterInput.trim();
      setFilterStrDebounced(text);
    }
  }

  onDestroy(() => {
    unsubscribeProcessStore();
  });
</script>

<Panel title="Logs" --panel-padding="0" --panel-overflow-y="hidden">
  <LocalLogProvider {workspace} {filterStr} {logs} let:events>
    <Logs {processIdToName} {events} />
  </LocalLogProvider>

  <div slot="bottom">
    <input type="text" placeholder="Filter..." bind:value={filterInput} />
  </div>
</Panel>

<style>
  input {
    width: 100%;
    border: none;
    border-top: 1px solid var(--layout-bg-color);
    background: var(--primary-bg-color);
    font-size: 16px;
    padding: 8px 12px;
    outline: none;
  }
</style>
