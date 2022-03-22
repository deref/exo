export const nonNull = <T>(x: T | null | undefined): T => {
  if (x == null) {
    throw new Error('expected non-null');
  }
  return x;
};

export type Optional<T> = T | null | undefined;
