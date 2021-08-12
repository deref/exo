export interface TaskDescription {
  id: string;
  jobId: string;
  parentId: string | null;
  name: string;
  status: Status;
  message: string;
  created: string;
  updated: string;
  started: string | null;
  finished: string | null;
  progress: TaskProgress | null;
}

export interface TaskProgress {
  current: number;
  total: number;
}

export type Status = 'pending' | 'running' | 'success' | 'failure';
