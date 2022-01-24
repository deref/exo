import { test, expect } from '@jest/globals';
import { parseSpans } from './parseSpans';
import type { Span } from './parseSpans';

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

  check('\u001b[31mRed\u001b[32mGreen\u001b[34mBlue\u001b[mReset', [
    {
      foreground: 'rgb(170, 0, 0)',
      text: 'Red',
    },
    {
      foreground: 'rgb(0, 170, 0)',
      text: 'Green',
    },
    {
      foreground: 'rgb(0, 0, 170)',
      text: 'Blue',
    },
    {
      text: 'Reset',
    },
  ]);
});
