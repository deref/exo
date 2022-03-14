<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import ProcessList from './ProcessList.svelte';
</script>

<Meta title="ProcessList" component={ProcessList} />

<Template let:args>
  <ProcessList {...args} />
</Template>

<Story
  name="Default"
  args={{
    workspace: {
      id: 'workspace1',
      displayName: 'example',
      processes: [],
    },
  }}
/>
