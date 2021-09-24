import { writable } from 'svelte/store';

export const isRunning = writable(true);
export const isAuthenticated = writable(true);
