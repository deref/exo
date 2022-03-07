<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import TaskTree from './TaskTree.svelte';
</script>

<Meta title="TaskTree" component={TaskTree} />

<Template let:args>
  <TaskTree {...args} />
</Template>

<Story
  name="Default"
  args={{
    value: {
      id: 'task001',
      name: 'Refresh',
      status: 'running',
      children: [
        {
          id: 'task002',
          status: 'success',
          name: 'Goes smoothly',
          progress: null,
          children: [],
        },
        {
          id: 'task003',
          status: 'failure',
          name: 'Goes poorly',
          progress: null,
          children: [],
        },
        {
          id: 'task004',
          status: 'pending',
          name: 'Will go',
          progress: null,
          children: [],
        },
      ],
    },
  }}
/>

<Story
  name="Progress"
  args={{
    value: {
      id: 'parent',
      name: 'Refresh',
      status: 'running',
      children: [0, 0.05, 0.1, 0.5, 1].map((progress, i) => ({
        id: `task-${i}`,
        status: progress == 1 ? 'success' : 'running',
        name: 'progress',
        progress,
        children: [],
      })),
    },
  }}
/>
