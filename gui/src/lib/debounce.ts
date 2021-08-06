type VoidFunc<TArgs extends any[]> = (...args: TArgs) => void;

export default function debounce<TArgs extends any[]>(
  f: VoidFunc<TArgs>,
  delayMs: number,
): VoidFunc<TArgs> {
  let timer: number | null = null;
  let lastArgs: TArgs | null = null;

  return (...args: TArgs) => {
    if (timer !== null) {
      window.clearTimeout(timer);
    }
    timer = window.setTimeout(() => {
      f(...(lastArgs as TArgs));
    }, delayMs);
    lastArgs = args;
  };
}
