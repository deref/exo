<script>
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import FormattedLogMessage from './FormattedLogMessage.svelte';
</script>

<Meta title="FormattedLogMessage" component={FormattedLogMessage} />

<Template let:args>
  <FormattedLogMessage {...args} />
</Template>

<Story
  name="Plain"
  args={{
    message: 'This is a plain message.',
  }}
/>

<Story
  name="With Link"
  args={{
    message:
      'This message uses a   string with a https://deref.io URL in it and formats  an anchor tag with the link.',
  }}
/>
