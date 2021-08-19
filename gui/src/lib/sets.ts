export const setsIdentical = <TItem>(...sets: Set<TItem>[]): boolean => {
  if (sets.length < 2) {
    return true;
  }

  const [first, second, ...rest] = sets;
  if (first.size !== second.size) {
    return false;
  }

  return (
    [...first.values()].every((item) => second.has(item)) &&
    setsIdentical(second, ...rest)
  );
};
