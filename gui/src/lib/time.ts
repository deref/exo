const padded = (int: number) => (int < 10 ? '0' + int : int);

export const shortTime = (timestamp: string): string => {
  const time = new Date(timestamp);

  return `${padded(time.getHours())}:${padded(time.getMinutes())}:${padded(
    time.getSeconds(),
  )}`;
};
