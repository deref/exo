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

export const logStyleFromHash = (s: string) => `--log-color:${hashPalette(s)[0]};--log-bg-color:${hashPalette(s)[2]};--log-border-color:${hashPalette(s)[1]}`


// $dark1: #32373c;
// $dark2: #3a4147;
// $dark3: #424f56;
// $grey-cold: #5f6b72;
// $grey1: #888888;
// $grey2: #aaaaaa;
// $grey3: #bbbbbb;
// $light1: #cccccc;
// $light2: #dddddd;
// $lightE: #eeeeee;
// $light3: #f9f9f9;
// $white: #ffffff;
// $red1: #db0000; // Color numbering: from darkest & richest to lightest & most pale.
// $red1-lighter: #f70e0e; // Color numbering: from darkest & richest to lightest & most pale.
// $red1-darker: #aa0000; // Color numbering: from darkest & richest to lightest & most pale.
// $red2: #ff8181;
// $red3: #f5d4d4;
// $orange1: #ff7008;
// $orange2: #ffb198;
// $orange3: #f5dcd4;
// $yellow1: #d38200; // For icons, text, etc
// $yellow2: #edc620; // For light backgrounds
// $yellow3: #fff2b3; // For very light backgrounds
// $green1: #00a508;
// $green1-darker: #008606;
// $green1-lighter: #00b909;
// $green2: #7dde75;
// $green3: #d6efd7;
// $blue1: ;
// $blue2: ;
// $blue3: ;
// $blue1-20: rgba(0, 113, 249, 0.2);
// $teal1: #119ab8;
// $teal2: #68c4d8;
// $teal3: #b0e5f0;
// $purple1: #8b00f9;
// $purple2: #bf8af3;
// $purple3: #eecaff;
