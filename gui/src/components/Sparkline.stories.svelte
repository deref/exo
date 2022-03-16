<script lang="ts">
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import Sparkline from './Sparkline.svelte';
</script>

<Meta title="Sparkline" component={Sparkline} />

<Template let:args>
  <Sparkline {...args} />
</Template>

<Story
  name="Default"
  args={{
    entries: [0, 1, 2, 3, 4, 3, 2, 1, 2, 3, 4],
    interactive: true,
  }}
/>
