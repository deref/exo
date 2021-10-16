// Based on https://svelte.dev/repl/acd92c9726634ec7b3d8f5f759824d15?version=3.43.0

interface Chord {
  alt?: boolean;
  shift?: boolean;
  command?: boolean; // Either control or meta.
  control?: boolean;
  meta?: boolean;
  code: string;
}

export interface ShortcutsParams {
  chords: Chord[];
  callback?(e: KeyboardEvent): void;
}

export const shortcuts = (node: HTMLElement, params: ShortcutsParams) => {
  let handler: (this: Window, e: KeyboardEvent) => any;
  const removeHandler = () => window.removeEventListener('keydown', handler),
    setHandler = () => {
      removeHandler();
      if (!params) {
        return;
      }
      handler = (e) => {
        const matched = params.chords.some((chord) => {
          if (
            !!chord.alt !== e.altKey ||
            !!chord.shift !== e.shiftKey ||
            chord.code !== e.code
          ) {
            return false;
          }
          if (!!chord.command) {
            return !!chord.command === (e.ctrlKey || e.metaKey);
          }
          return !!chord.control === e.ctrlKey && !!chord.meta === e.metaKey;
        });
        if (!matched) {
          return;
        }
        e.preventDefault();
        params.callback ? params.callback(e) : node.click();
      };
      window.addEventListener('keydown', handler);
    };
  setHandler();
  return {
    update: setHandler,
    destroy: removeHandler,
  };
};

export interface ShortcutParams extends Chord {
  callback?(e: KeyboardEvent): void;
}

export const shortcut = (
  node: HTMLElement,
  params: ShortcutParams | undefined,
) => {
  if (!params) {
    return;
  }
  const { callback, ...chord } = params;
  shortcuts(node, { chords: [chord], callback });
};
