// Based on https://svelte.dev/repl/acd92c9726634ec7b3d8f5f759824d15?version=3.43.0

interface Params {
  alt?: boolean;
  shift?: boolean;
  control?: boolean;
  code?: string;
  callback?(e: KeyboardEvent): void;
}

export const shortcut = (node: HTMLElement, params: Params) => {
  let handler: (this: Window, e: KeyboardEvent) => any;
  const removeHandler = () => window.removeEventListener('keydown', handler),
    setHandler = () => {
      removeHandler();
      if (!params) return;
      handler = (e) => {
        if (
          !!params.alt != e.altKey ||
          !!params.shift != e.shiftKey ||
          !!params.control != (e.ctrlKey || e.metaKey) ||
          params.code != e.code
        )
          return;
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
