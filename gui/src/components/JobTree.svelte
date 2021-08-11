<script lang="ts">
  import type { TaskNode } from './TaskTree.svelte';
  import TaskTree from './TaskTree.svelte';

  import type { TaskDescription } from '../lib/tasks/types';

  export let jobId: string;
  export let tasks: TaskDescription[];

  const makeTaskTree = (
    jobId: string,
    tasks: TaskDescription[],
  ): TaskNode | null => {
    const nodes = new Map<string, TaskNode>();
    const getNode = (id: string): TaskNode => {
      let node = nodes.get(id);
      if (node == null) {
        node = {
          id,
          name: undefined as any,
          status: undefined as any,
          progress: null,
          children: [],
        };
        nodes.set(id, node);
      }
      return node;
    };

    for (const task of tasks) {
      if (task.jobId !== jobId) {
        continue;
      }
      const child = getNode(task.id);
      child.name = task.name;
      child.status = task.status;
      if (task.progress !== null) {
        child.progress = task.progress.current / task.progress.total;
      }

      if (task.parentId !== null) {
        const parent = getNode(task.parentId);
        parent.children.push(child);
      }
    }
    return getNode(jobId);
  };

  $: tree = makeTaskTree(jobId, tasks);
</script>

{#if tree}
  <TaskTree value={tree} />
{/if}
