import { test, expect } from '@jest/globals';
import { parseSpans, Span } from './parseSpans';

test('parseSpans', () => {
  const check = (input: string, expected: Span[]) => {
    expect(parseSpans(input)).toEqual(expected);
  };

  check('', []);
  check('asdf', [
    {
      text: 'asdf',
    },
  ]);
  check('https://foo.com', [
    {
      href: 'https://foo.com',
      text: 'https://foo.com',
    },
  ]);
  check('foo https://foo.com bar', [
    {
      text: 'foo ',
    },
    {
      href: 'https://foo.com',
      text: 'https://foo.com',
    },
    {
      text: ' bar',
    },
  ]);
  check('foo <https://foo.com> bar', [
    {
      text: 'foo <',
    },
    {
      href: 'https://foo.com',
      text: 'https://foo.com',
    },
    {
      text: '> bar',
    },
  ]);

  check('\u001b[31mRed\u001b[32mGreen\u001b34mBlue', [
    {
      forground: '#ff0000',
      text: 'Red',
    },
    {
      forground: '#00ff00',
      text: 'Green',
    },
    {
      forground: '#0000ff',
      text: 'Blue',
    },
  ]);
});
