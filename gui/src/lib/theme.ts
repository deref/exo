import { writable } from 'svelte/store';

export const themeLocalStorageKey = 'io.deref.exo/gui-theme';

export const themeOptions = ['auto', 'light', 'dark', 'black'];

type ThemeName = 'auto' | 'light' | 'dark' | 'black';

function createTheme() {
  const ls = localStorage.getItem(themeLocalStorageKey) || 'auto';

  const { subscribe, set } = writable<ThemeName>(<ThemeName>ls);

  const localSyncedSet = (t: ThemeName) => {
    localStorage.setItem(themeLocalStorageKey, t);
    set(t);
  };

  return {
    subscribe,
    apply: (t: ThemeName) => localSyncedSet(t),
  };
}

export const theme = createTheme();
