<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import LogFilterBar from './LogFilterBar.svelte';
</script>

<Meta title="LogFilterBar" component={LogFilterBar} />

<Template let:args>
  <LogFilterBar {...args} />
</Template>

<Story
  name="Default"
  args={{
    filterString: '',
  }}
/>
