<script lang="ts">
  import { parseSpans } from './parseSpans';
  import type { Span } from './parseSpans';

  export let message: string = '';

  const spans = parseSpans(message);

  const spanProps = (span: Span): Record<string, string | boolean> => {
    const styles: string[] = [];
    const classes: string[] = [];
    const addStyle = (key: string, value: string) => {
      styles.push(`${key}: ${value}`);
    };
    if (span.foreground != null) {
      addStyle('color', span.foreground);
    }
    if (span.background != null) {
      addStyle('background', span.background);
    }
    switch (span.style) {
      case 'bold':
        addStyle('font-weight', 'bold');
        break;
      case 'faint':
        addStyle('font-weight', 'lighter');
        break;
      case 'italic':
        addStyle('font-style', 'italic');
        break;
      case 'underline':
        addStyle('text-decoration', 'underline');
        break;
      case 'blink':
        classes.push('exo-message-blink');
        break;
      case 'invert':
        classes.push('exo-message-invert');
        break;
      case 'strike':
        addStyle('text-decoration', 'line-through');
        break;
    }
    const props: Record<string, string | boolean> = {};
    if (styles.length > 0) {
      props['style'] = styles.join(':');
    }
    if (classes.length > 0) {
      props['class'] = classes.join(' ');
    }
    return props;
  };
</script>

<span>
  {#each spans as span}
    {#if span.href != null}
      <a href={span.href} {...spanProps(span)}>{span.text}</a>
    {:else}
      <span {...spanProps(span)}>{span.text}</span>
    {/if}
  {/each}
</span>

<style>
  span {
    font-family: var(--font-mono);
    font-variant-ligatures: var(--preferred-ligatures-logs);
    white-space: pre-wrap;
  }

  :global(.exo-message-blink) {
    animation: blink 1s steps(1) infinite;
  }

  @keyframes blink {
    50% {
      color: transparent;
    }
  }

  :global(.exo-message-invert) {
    filter: invert(1);
  }
</style>
