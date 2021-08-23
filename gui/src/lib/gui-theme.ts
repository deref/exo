import { writable } from 'svelte/store';

export const guiThemeLocalStorageKey = 'io.deref.exo/gui-theme';

function createGuiTheme() {
  const ls = localStorage.getItem(guiThemeLocalStorageKey) || 'auto';

  const { subscribe, set } = writable(ls);

  const localSyncedSet = (s: string) => {
    set(s);
    localStorage.setItem(guiThemeLocalStorageKey, s);
  };

  return {
    subscribe,
    auto: () => localSyncedSet('auto'),
    light: () => localSyncedSet('light'),
    dark: () => localSyncedSet('dark'),
    black: () => localSyncedSet('black'),
  };
}

export const guiTheme = createGuiTheme();

export const setGuiTheme = (theme: string) => {
  for (const option of ['auto', 'light', 'dark', 'black']) {
    if (theme === option) {
      document.body.classList.add(option);
    } else {
      document.body.classList.remove(option);
    }
  }
};
