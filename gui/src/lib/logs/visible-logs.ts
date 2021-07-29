import { processes } from '../process/store';
import { derived, writable } from 'svelte/store';

const visibleLogsKey = 'exo:logs:visible';

// Store managing the list of visible log providers,
// XXX: This uses a global localStorage instance, but it should probably be persisted to the server
// and scoped to the workspace. The next refactor should move this to a context-providing component
// that is initialized whenever a workspace is loaded.
const initialHiddenLogs = new Set(
  JSON.parse(localStorage.getItem(visibleLogsKey) ?? '[]') as string[],
); // TODO: Validate w/ runtype.
export const hiddenLogsStore = writable<Set<string>>(initialHiddenLogs);

export const visibleLogsStore = derived(
  [processes, hiddenLogsStore],
  ([processes, hiddenLogs]) => {
    const set = new Set<string>();
    if ('data' in processes) {
      for (const process of processes.data) {
        set.add(process.id);
      }
      for (const id of hiddenLogs) {
        set.delete(id);
      }
    }
    return set;
  },
);

export const setLogVisibility = (processId: string, visible: boolean) => {
  hiddenLogsStore.update((hiddenLogs) => {
    if (visible) {
      hiddenLogs.delete(processId);
    } else {
      hiddenLogs.add(processId);
    }

    localStorage.setItem(
      visibleLogsKey,
      JSON.stringify([...hiddenLogs.values()]),
    );
    return hiddenLogs;
  });
};
