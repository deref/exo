<script>
  import { Meta, Template, Story } from '@storybook/addon-svelte-csf';
  import { range } from '../../lib/arrays';
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

<Story
  name="Colors: 4 bit"
  args={{
    message: [
      `\u001b[30mX\u001b[31mX\u001b[32mX\u001b[33mX\u001b[34mX\u001b[35mX\u001b[36mX\u001b[37mX\u001b[0m `,
      `\u001b[40mX\u001b[41mX\u001b[42mX\u001b[43mX\u001b[44mX\u001b[45mX\u001b[46mX\u001b[47mX\u001b[0m `,
      `\u001b[90mX\u001b[91mX\u001b[92mX\u001b[93mX\u001b[94mX\u001b[95mX\u001b[96mX\u001b[97mX\u001b[0m `,
      `\u001b[100mX\u001b[101mX\u001b[102mX\u001b[103mX\u001b[104mX\u001b[105mX\u001b[106mX\u001b[107mX\u001b[0m `,
    ].join(''),
  }}
/>

<Story
  name="Colors: 8 bit"
  args={{
    message:
      range(256)
        .map((i) => `\u001b[38;5;${i};m ${i}`)
        .join('') +
      range(256)
        .map((i) => `\u001b[48;5;${i};m ${i}`)
        .join(''),
  }}
/>

<Story
  name="Text Decorations"
  args={{
    message: `\u001b[1mBold\u001b[0m \u001b[2mFaint\u001b[0m \u001b[3mItalic\u001b[0m \u001b[4mUnderline\u001b[0m \u001b[9mStrikethrough`,
  }}
/>
