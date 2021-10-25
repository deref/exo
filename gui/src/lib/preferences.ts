import { writable } from 'svelte/store';

const key = 'io.deref.exo/preferences';

type Preferences = Record<string, string>;

let defaults = {
  'main-font-size': '16px',
  'log-font-size': '15px',
  'ligatures-logs': 'none',
  'ligatures-code': 'none',
};

function createStoredPreferences() {
  let obj: Preferences = { ...defaults };

  for (const variable of Object.keys(defaults)) {
    const localItem = localStorage.getItem(key + '/' + variable);

    if (localItem !== null) {
      obj[variable] = localItem;
    }
  }

  const { subscribe, set } = writable<Preferences>(<Preferences>obj);

  const localSyncedSet = (prefs: Preferences) => {
    for (const variable of Object.keys(defaults)) {
      localStorage.setItem(key + '/' + variable, prefs[variable]);
    }

    set(prefs);
  };

  return {
    subscribe,
    apply: (prefs: Preferences) => localSyncedSet(prefs),
    reset: () => localSyncedSet({ ...defaults }),
  };
}

export const preferences = createStoredPreferences();
