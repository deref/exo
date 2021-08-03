import type { SvelteComponentTyped } from 'svelte';

export type SvelteConstructor<T> = new (
  ...args: any
) => SvelteComponentTyped<T>;
