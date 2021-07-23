// This is a palette of colors intended to be different enough from one another
// that distinct processes look distinct in the log viewer.
export const palette = [
  [
  '#0071f9',
  '#8abaf3',
  '#bbd7f9',
  ],
  [
  '#008606',
  '#7dde75',
  '#d6efd7',
  ],
  [
  '#ff7008',
  '#ffb198',
  '#f5dcd4',
  ],
  [
  '#8b00f9',
  '#bf8af3',
  '#eecaff',
  ],
  [
  '#d38200',
  '#edc620',
  '#fff2b3',
  ],
  [
  '#db0000',
  '#ff8181',
  '#f5d4d4',
  ],
]

// This computes an integer hash which can be used to select a palette.
export const hashString = (s: string) => s.split('').map(c => c.charCodeAt(0)).reduce((a,c) => a + c)

// This selects a color palette based on the hash of a string.
export const hashPalette = (s: string) => palette[hashString(s) % palette.length]

// This computes a HTML style attribute string for colored logs.
export const logStyleFromHash = (s: string) => {
  // Format a single color variable.
  const f = (varName, paletteIndex) => { 
    return `--${varName}:${hashPalette(s)[paletteIndex]};`
  }
  // Combine styles into one string.
  return f("log-color",0) + f("log-border-color",1) + f("log-bg-color",2)
}
