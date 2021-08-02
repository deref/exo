<script lang="ts">
  import * as monaco from 'monaco-editor';
  import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';

  import { onMount } from 'svelte';

  export let id: string | undefined;
  export let value: string = '';
  export let language: string;
  export let height: string = '200px';

  (self as any).MonacoEnvironment = {
    getWorker(_: unknown, label: string) {
      return new editorWorker();
    },
  };

  let container: HTMLDivElement | null = null;

  onMount(() => {
    const editor = monaco.editor.create(container!, {
      value,
      language,
      lineNumbers: 'off',
      minimap: {
        enabled: false,
      },
    });
    editor.onDidChangeModelContent((event) => {
      value = editor.getValue();
    });

    return () => {
      editor.dispose();
    };
  });
</script>

<div {id} bind:this={container} style={`height: ${height}`} />
