<script lang="ts">
  export let code: string | undefined = undefined;
</script>

<pre><code>{code}</code></pre>

<style>
  pre {
    background: var(--code-bg-color);
    box-shadow: var(--heavy-3d-box-shadow);
    font-size: 15px;
    font-weight: 400;
    border-radius: 0.33em;
    padding: 0.667em 1em;
    overflow-x: auto;
  }
</style>
