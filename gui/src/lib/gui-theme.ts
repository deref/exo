import { writable } from 'svelte/store';

export const themeLocalStorageKey = 'io.deref.exo/gui-theme';

type ThemeName = 'auto' | 'light' | 'dark' | 'black';

export const applyBodyTheme = (theme: ThemeName) => {
  for (const option of ['auto', 'light', 'dark', 'black']) {
    document.body.classList.toggle(option, theme === option);
  }
};

function createTheme() {
  const ls = localStorage.getItem(themeLocalStorageKey) || 'auto';

  const { subscribe, set } = writable<ThemeName>(<ThemeName>ls);

  const localSyncedSet = (t: ThemeName) => {
    set(t);
    localStorage.setItem(themeLocalStorageKey, t);
  };

  return {
    subscribe,
    apply: (t: ThemeName) => {
      localSyncedSet(t);
      applyBodyTheme(t);
    },
  };
}

export const theme = createTheme();
