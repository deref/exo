export interface Span {
  href?: string;
  text: string;
  foreground?: string;
  background?: string;
  style?: Style;
}

export type Style =
  | 'bold'
  | 'faint'
  | 'italic'
  | 'underline'
  | 'blink'
  | 'invert'
  | 'strike';

const isLinkOpenBoundary = (c: string) => {
  switch (c) {
    case '':
    case ' ':
    case '\n':
    case '<':
      return true;
    default:
      return false;
  }
};

const isLinkCloseBoundary = (c: string) => {
  switch (c) {
    case '':
    case ' ':
    case '\n':
    case '<':
      return true;
    default:
      return false;
  }
};

const standardColors: string[] = [
  'rgb(0, 0, 0)', // Black
  'rgb(170, 0, 0)', // Red
  'rgb(0, 170, 0)', // Green
  'rgb(170, 85, 0)', // Yellow
  'rgb(0, 0, 170)', // Blue
  'rgb(170, 0, 170)', // Magenta
  'rgb(0, 170, 170)', // Cyan
  'rgb(170, 170, 170)', // White
];

const brightColors: string[] = [
  'rgb(85, 85, 85)', // Gray
  'rgb(255, 85, 85)', // Red
  'rgb(85, 255, 85)', // Green
  'rgb(255, 255, 85)', // Yellow
  'rgb(85, 85, 255)', // Blue
  'rgb(255, 85, 255)', // Magenta
  'rgb(85, 255, 255)', // Cyan
  'rgb(255, 255, 255)', // White
];

const lookupColor256 = (i: number): string => {
  if (i < 0) {
    return '#000000';
  }
  if (i < 8) {
    return standardColors[i];
  }
  if (i < 16) {
    return brightColors[i - 8];
  }
  if (i < 232) {
    // Color cube of size 6.
    // TODO: Does this produce the correct colors?
    const x = i - 16;
    const r = ((x >> 4) % 6) * 51;
    const g = ((x >> 2) % 6) * 51;
    const b = ((x >> 0) % 6) * 51;
    return `rgb(${r},${g},${b})`;
  }
  if (i < 256) {
    // Grayscale.
    const x = ((i - 232) / 24) * 256;
    return `rgb(${x}, ${x}, ${x})`;
  }
  return '000000';
};

const lookupColor24bit = (codes: number[]): string => {
  const r = 200;
  const g = 200;
  const b = 200;
  return `rgb(${r}, ${g}, ${b})`;
};

const lookupColor = (codes: number[]) => {
  switch (codes[1]) {
    case 5:
      return lookupColor256(codes[2]);
    case 2:
      return lookupColor24bit(codes);
    default:
      return '000000';
  }
};

export const parseSpans = (input: string): Span[] => {
  // TODO: Create reusable state machine so that multiple independent lines
  // can have persistent color styles across lines.
  let foreground: string | null = null;
  let background: string | null = null;
  let style: Style | null = null;
  let inLink = false;
  let href = '';
  let text = '';

  const spans: Span[] = [];

  const emitSpan = () => {
    if (text.length > 0) {
      const span: Span = {
        text,
      };
      const setOptional = <TKey extends keyof Span>(
        key: TKey,
        value: Span[TKey] | null,
      ) => {
        if (value) {
          span[key] = value;
        }
      };
      setOptional('href', href);
      setOptional('foreground', foreground);
      setOptional('background', background);
      setOptional('style', style);
      spans.push(span);
      text = '';
    }
  };

  let pos = 0;
  let prev = '';
  while (pos <= input.length) {
    const cur = pos === input.length ? '' : input[pos];
    const rest = input.slice(pos);

    // End the current link?
    if (inLink) {
      const endOfLink = isLinkCloseBoundary(cur);
      if (endOfLink) {
        href = text;
        emitSpan();
        href = '';
        inLink = false;
      }
    }

    // Start a new link?
    if (!inLink && isLinkOpenBoundary(prev) && /^https?:\/\//i.test(rest)) {
      emitSpan();
      inLink = true;
    }

    // Handle control sequences.
    const control = rest.match(/^\u001b\[([\d;]+)m/);
    if (control) {
      emitSpan();
      let codes = control[1].split(';').map((n) => parseInt(n, 10));
      if (codes.length === 0) {
        codes = [0];
      }
      const command = codes[0];
      switch (command) {
        case 0: {
          // Reset.
          foreground = null;
          background = null;
          style = null;
          break;
        }

        case 23: {
          if (style === 'bold' || style === 'italic') {
            style = null;
          }
          break;
        }

        case 1:
          style = 'bold';
          break;
        case 2:
          style = 'faint';
          break;
        case 3:
          style = 'italic';
          break;
        case 4:
          style = 'underline';
          break;
        case 5:
        case 6:
          style = 'blink';
          break;
        case 7:
          style = 'invert';
          break;
        case 9:
          style = 'strike';
          break;

        case 38:
          foreground = lookupColor(codes);
          break;
        case 48:
          background = lookupColor(codes);
          break;

        default: {
          if (30 <= command && command <= 37) {
            foreground = standardColors[command - 30];
          } else if (40 <= command && command <= 47) {
            background = standardColors[command - 40];
          } else if (90 <= command && command <= 97) {
            foreground = brightColors[command - 90];
          } else if (100 <= command && command <= 107) {
            background = brightColors[command - 100];
          }
        }
      }
      pos += control[0].length;
    } else {
      text += cur;
      pos += 1;
    }

    prev = cur;
  }

  // Flush.
  emitSpan();

  return spans;
};
