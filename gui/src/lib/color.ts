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

export const hashString = (s: string) => s.split('').map(c => c.charCodeAt(0)).reduce((a,c) => a + c)

export const hashPalette = (s: string) => palette[hashString(s) % palette.length]

export const logStyleFromHash = (s: string) => {

  const f = (varName, paletteIndex) => { 
    
    return `--${varName}:${hashPalette(s)[paletteIndex]};`
  }

  return f("log-color",0) + f("log-border-color",1) + f("log-bg-color",2)

}
