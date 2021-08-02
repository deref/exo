import { test, expect } from '@jest/globals';
import { parseSpans, Span } from './parseSpans';

test('parseSpans', () => {
  const check = (input: string, expected: Span[]) => {
    expect(parseSpans(input)).toEqual(expected);
  };

  check('', []);
  check('asdf', [
    {
      type: 'plain',
      text: 'asdf',
    },
  ]);
  check('https://foo.com', [
    {
      type: 'link',
      href: 'https://foo.com',
      text: 'https://foo.com',
    },
  ]);
  check('foo https://foo.com bar', [
    {
      type: 'plain',
      text: 'foo ',
    },
    {
      type: 'link',
      href: 'https://foo.com',
      text: 'https://foo.com',
    },
    {
      type: 'plain',
      text: ' bar',
    },
  ]);
  check('foo <https://foo.com> bar', [
    {
      type: 'plain',
      text: 'foo <',
    },
    {
      type: 'link',
      href: 'https://foo.com',
      text: 'https://foo.com',
    },
    {
      type: 'plain',
      text: '> bar',
    },
  ]);
});
