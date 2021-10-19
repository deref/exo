import { getContext, setContext } from 'svelte';

export const contextKey = {};

export interface Menu {
  close(): void;
}

export const setMenuContext = (menu: Menu) => {
  setContext(contextKey, menu);
};

export const getMenuContext = (): Menu | undefined =>
  getContext(contextKey) as Menu | undefined;
