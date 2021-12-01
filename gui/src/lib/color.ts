// This computes an integer hash which can be used to select a palette.
export const hashString = (s: string) =>
  s
    .split('')
    .map((c) => c.charCodeAt(0))
    .reduce((a, c) => a + c, 0);

export const hashDegree = (s: string) =>
  Math.round(hashString(s) * Math.PI * 100) % 360;

export const textColor = (deg: number) => `hsl(${deg}, 95%, 30%)`;
export const textHoverColor = (deg: number) => `hsl(${deg}, 95%, 20%)`;
export const bgColor = (deg: number) => `hsl(${deg}, 65%, 92%)`;
export const bgHoverColor = (deg: number) => `hsl(${deg}, 65%, 87%)`;

export const darkTextColor = (deg: number) => `hsl(${deg}, 75%, 75%)`;
export const darkTextHoverColor = (deg: number) => `hsl(${deg}, 75%, 90%)`;
export const darkBgColor = (deg: number) => `hsl(${deg}, 80%, 6%)`;
export const darkBgHoverColor = (deg: number) => `hsl(${deg}, 80%, 12%)`;

// This computes a HTML style attribute string for colored logs.
export const logStyleFromHash = (s: string) => {
  const d = hashDegree(s);
  // Combine styles into one string.
  return `--light-log-color:${textColor(d)};
  --light-log-bg-color:${bgColor(d)};
  --light-log-hover-color:${textHoverColor(d)};
  --light-log-bg-hover-color:${bgHoverColor(d)};
  --dark-log-color:${darkTextColor(d)};
  --dark-log-bg-color:${darkBgColor(d)};
  --dark-log-hover-color:${darkTextHoverColor(d)};
  --dark-log-bg-hover-color:${darkBgHoverColor(d)}`;
};
