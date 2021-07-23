import { writable } from 'svelte/store';

const visibleLogsKey = 'exo:logs:visible';

// Store managing the list of visible log providers,
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
