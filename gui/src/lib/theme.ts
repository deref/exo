import { writable } from 'svelte/store';

const themeLocalStorageKey = 'io.deref.exo/gui-theme';

export const themeOptions = ['auto', 'light', 'dark', 'black'] as const;

type ThemeOption = typeof themeOptions[number];

function createTheme() {
  const ls = localStorage.getItem(themeLocalStorageKey) || 'auto';

  const { subscribe, set } = writable<ThemeOption>(<ThemeOption>ls);

  const localSyncedSet = (t: ThemeOption) => {
    localStorage.setItem(themeLocalStorageKey, t);
    set(t);
  };

  return {
    subscribe,
    apply: (t: ThemeOption) => localSyncedSet(t),
  };
}

export const theme = createTheme();
