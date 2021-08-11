<script lang="ts">
  import { parseSpans, Span } from './parseSpans';

  export let message: string = '';

  const spans = parseSpans(message);

  const spanStyle = (span: Span): string => {
    const styles: string[] = [];
    const addStyle = (key: string, value: string) => {
      styles.push(`${key}: ${value}`);
    };
    if (span.forground != null) {
      addStyle('color', span.forground);
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
      case 'blink':
      case 'invert':
        // Not supported (yet?).
        break;
      case 'strike':
        addStyle('text-decoration', 'line-through');
        break;
    }
    return styles.join(';');
  };
</script>

<span>
  {#each spans as span}
    {#if span.href != null}
      <a href={span.href} style={spanStyle(span)}>{span.text}</a>
    {:else}
      <span style={spanStyle(span)}>{span.text}</span>
    {/if}
  {/each}
</span>

<style>
  a,
  span {
    white-space: pre;
  }
</style>
