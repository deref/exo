import type { ProcessSpec } from "src/lib/api";
import * as shell from 'shell-quote';

const shellQuote = (s: string) => shell.quote([s]);

export interface ParseResult {
  spec: Partial<ProcessSpec>; 
  error: Error | null;
}

export const parseScript = (script: string): ParseResult => {
  const lines = script.split('\n').map(line => line.trim());
  
  const spec: ProcessSpec = {
    program: '',
    arguments: [],
  }
  let error: Error | null = null;

  let lineIndex = 0;
lines:
  for (const line of lines) {
    lineIndex++;

    const setError = (message: string) => {
      if (error) {
        return;
      }
      error = new Error(`line ${lineIndex}: ${message}`);
    };

    const entries = shell.parse(line);
    if (entries.length < 1) {
      continue lines;
    }
    const words: string[] = [];
  entries:
    for (const entry of entries) {
      if (typeof entry === 'string') {
        words.push(entry);
      } else if ('comment' in entry) {
        continue entries;
      } else if (entry.op === 'glob') {
        setError('unexpected glob');
        continue lines;
      } else {
        setError(`unexpected operator: ${entry.op}`)
      }
    }
    const [program, ...args] = words;
    switch (program) {
      case 'cd': {
        if (args.length !== 1) {
          setError('cd expects only one argument');
        }
        spec.directory = args[0];
        break;
      }
      case 'export': {
        if (args.length !== 1) {
          setError('export expects only one argument of the form NAME=VALUE');
          break;
        }
        const [kvp] = args;
        const parts = kvp.split('=');
        if (parts.length !== 2) {
          setError('export expects only one argument of the form NAME=VALUE');
          break;
        }
        const [key, val] = parts;
        spec.environment = spec.environment ?? {};
        spec.environment[key] = val;
        break;
      }
      default: {
        if (spec.program) {
          setError('script must end with a single program invocation');
          break;
        }
        spec.program = program;
        spec.arguments = args;
        break;
      }
    }
  }
  
  return { spec, error };
};

export const generateScript = (spec: ProcessSpec): string => {
  let script = '';
  if (spec.directory) {
    script += `cd ${shellQuote(spec.directory)}\n`;
  }
  for (const [k, v] of Object.entries(spec.environment ?? {})) {
    script += `export ${shellQuote(k)}=${shellQuote(v)}\n`;
  }
  const words: string[] = [];
  if (spec.program !== '') {
    words.push(spec.program);
  }
  words.push(...spec.arguments);
  script += words.map(word => shellQuote(word)).join(' ');
  return script.trim();
}
