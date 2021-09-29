<script lang="ts">
  import Panel from './Panel.svelte';
  import { shortcuts } from '../lib/actions/shortcut';
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
  const initProcessIdToName: Record<string, string> = {
    [workspace.id]: 'EXO',
  };
  let processIdToName = initProcessIdToName;
  const unsubscribeProcessStore = processes.subscribe((processData) => {
    if (processData.stage === 'success' || processData.stage === 'refetching') {
      processIdToName = processData.data.reduce(
        (acc, processDescription) => ({
          ...acc,
          [processDescription.id]: processDescription.name,
        }),
        initProcessIdToName,
      );
    }
  });

  const getComponentName = (id: string) => processIdToName[id] ?? null;

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

<LocalLogProvider
  {workspace}
  {filterStr}
  streams={[workspace.id, ...logs]}
  let:events
  let:clearEvents
>
  <Panel title="Logs" --panel-padding="0" --panel-overflow-y="hidden">
    <Logs {getComponentName} {events} />
    <div class="bottom" slot="bottom">
      <input type="text" placeholder="Filter..." bind:value={filterInput} />
      <button
        use:shortcuts={{
          chords: [
            {
              meta: true,
              code: 'KeyK',
            },
            {
              control: true,
              code: 'KeyL',
            },
          ],
        }}
        on:click={(e) => {
          clearEvents();
        }}
      >
        Clear Logs
      </button>
    </div>
  </Panel>
</LocalLogProvider>

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

  .bottom {
    display: flex;
  }

  button {
    font-size: 80%;
    width: 100px;
  }
</style>
