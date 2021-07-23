import { writable } from 'svelte/store';

const visibleLogsKey = 'exo:logs:visible';

// Store managing the list of visible log providers,
// XXX: This uses a global localStorage instance, but it should probably be persisted to the server
// and scoped to the workspace. The next refactor should move this to a context-providing component
// that is initialized whenever a workspace is loaded.
const initialVisibleLogs = new Set(JSON.parse(localStorage.getItem(visibleLogsKey) ?? '[]') as string[]); // TODO: Validate w/ runtype.
export const visibleLogsStore = writable<Set<string>>(initialVisibleLogs);

export const toggleLogVisibility = (processId: string) => {
    visibleLogsStore.update(visibleLogs => {
        if (visibleLogs.has(processId)) {
            visibleLogs.delete(processId);
        } else {
            visibleLogs.add(processId);
        }

        localStorage.setItem(visibleLogsKey, JSON.stringify([...visibleLogs.values()]));
        return visibleLogs;
    })
};
