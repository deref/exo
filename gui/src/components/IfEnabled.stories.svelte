<script>
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import IfEnabled from './IfEnabled.svelte';
</script>

<Meta title="IfEnabled" component={IfEnabled} />

<Template let:args>
  <IfEnabled feature={args.feature}>
    {args.content}
  </IfEnabled>
</Template>

<Story
  name="Enabled"
  args={{
    feature: 'sandbox test',
    content: 'Enabled feature: sandbox test',
  }}
/>

<Story
  name="Disabled"
  args={{
    feature: 'bogus feature',
    content: 'Disabled feature: bogus feature',
  }}
/>
