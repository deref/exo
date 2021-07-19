import { writable } from 'svelte/store';

export const workspaceId = writable<string | null>(null);
