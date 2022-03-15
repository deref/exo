<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import ComponentStack from './ComponentStack.svelte';
  import type { Component } from './ComponentStack.svelte';

  const componentDefaults: Partial<Component> = {
    logsVisible: true,
    reconciling: false,
    running: true,
    url: '/#',
    editUrl: '/#edit',
  };
</script>

<Meta title="ComponentStack" component={ComponentStack} />

<Template let:args>
  <ComponentStack {...args} />
</Template>

<Story
  name="Default"
  args={{
    components: [
      {
        id: 'backend',
        name: 'Backend',
        ...componentDefaults,
      },
      {
        id: 'frontend',
        name: 'Frontend',
        ...componentDefaults,
        running: false,
      },
      {
        id: 'helper',
        name: 'helper',
        ...componentDefaults,
        reconciling: true,
        logsVisible: false,
      },
    ],
  }}
/>
