import { shortTime } from '../../lib/time';
import { logStyleFromHash } from '../../lib/color';
import type { LogEvent as ApiLogEvent } from 'src/lib/logs/types';
import type { LogEvent as UILogEvent } from './types';

const friendlyName = (
  log: string,
  processIdToName: Record<string, string>,
): string => {
  // SEE NOTE [LOG_COMPONENTS].
  const [procId, stream] = log.split(':');
  const procName = processIdToName[procId];
  return procName ? (stream ? `${procName}:${stream}` : procName) : log;
};

export const formatLogs = (
  events: ApiLogEvent[],
  processIdToName: Record<string, string>,
): UILogEvent[] =>
  events.map((event) => ({
    id: event.id,
    style: logStyleFromHash(event.stream),
    time: {
      short: shortTime(event.timestamp),
      full: event.timestamp,
    },
    name: friendlyName(event.stream, processIdToName),
    message: event.message,
  }));
