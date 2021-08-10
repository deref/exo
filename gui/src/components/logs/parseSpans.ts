export interface Span {
  href?: string;
  text: string;
}

export const parseSpans = (input: string): Span[] => {
  const spans: Span[] = [];
  //const re = /(?<=^\s\<)https?:\/\/(.+?)(?=$\s\>)/gi;
  const re = /(^|(?<=[\s<]))https?:\/\/(.+?)($|(?=[\s>]))/gi;
  let left = 0;
  while (left < input.length) {
    const match = re.exec(input);
    const right = match ? re.lastIndex - match[0].length : input.length;
    // If the search skipped some plain text or reached end, capture it.
    if (left < right) {
      spans.push({
        text: input.substring(left, right),
      });
      left = right;
    }
    if (match) {
      const text = match[0];
      spans.push({
        href: text,
        text,
      });
      left += text.length;
    }
  }
  return spans;
};
