// This computes an integer hash which can be used to select a palette.
export const hashString = (s: string) =>
  s
    .split('')
    .map((c) => c.charCodeAt(0))
    .reduce((a, c) => a + c);

export const hashDegree = (s: string) =>
  Math.round(hashString(s) * Math.PI * 100) % 360;

export const textColor = (deg: number) => `hsl(${deg}, 95%, 30%)`;
export const borderColor = (deg: number) => `hsl(${deg}, 65%, 75%)`;
export const bgColor = (deg: number) => `hsl(${deg}, 65%, 90%)`;

// This computes a HTML style attribute string for colored logs.
export const logStyleFromHash = (s: string) => {
  const d = hashDegree(s);
  // Combine styles into one string.
  return `--log-color:${textColor(d)};--log-border-color:${borderColor(
    d,
  )};--log-bg-color:${bgColor(d)}`;
};
