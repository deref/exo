// This computes an integer hash which can be used to select a palette.
export const hashString = (s: string) =>
  s
    .split('')
    .map((c) => c.charCodeAt(0))
    .reduce((a, c) => a + c);

export const hashDegree = (s: string) =>
  Math.round(hashString(s) * Math.PI * 100) % 360;

export const textColor = (deg: number) => `hsl(${deg}, 95%, 30%)`;
export const textHoverColor = (deg: number) => `hsl(${deg}, 95%, 20%)`;
export const bgColor = (deg: number) => `hsl(${deg}, 65%, 92%)`;
export const bgHoverColor = (deg: number) => `hsl(${deg}, 65%, 87%)`;

// This computes a HTML style attribute string for colored logs.
export const logStyleFromHash = (s: string) => {
  const d = hashDegree(s);
  // Combine styles into one string.
  return `--log-color:${textColor(d)};--log-bg-color:${bgColor(
    d,
  )};--log-hover-color:${textHoverColor(d)};--log-bg-hover-color:${bgHoverColor(
    d,
  )}`;
};
