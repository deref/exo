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
          children: [],
        },
        {
          id: 'task003',
          status: 'failure',
          name: 'Goes poorly',
          children: [],
        },
        {
          id: 'task003',
          status: 'pending',
          name: 'Will go',
          children: [],
        },
      ],
    },
  }}
/>
