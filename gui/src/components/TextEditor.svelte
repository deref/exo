<script lang="ts">
  import * as monaco from 'monaco-editor';
  import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
  import { exo_dark, exo_light } from '../lib/monaco-theme';

  import { onMount } from 'svelte';

  export let id: string | undefined;
  export let value: string = '';
  export let language: string;
  export let height: string = '240px';

  (self as any).MonacoEnvironment = {
    getWorker(_: unknown, label: string) {
      return new editorWorker();
    },
  };

  let container: HTMLDivElement | null = null;

  onMount(() => {
    monaco.editor.defineTheme('exo-light', exo_light);
    monaco.editor.defineTheme('exo-dark', exo_dark);

    const setTheme = (dark: boolean) => {
      if (dark) {
        monaco.editor.setTheme('exo-dark');
      } else {
        monaco.editor.setTheme('exo-light');
      }
    };

    const mqList = window.matchMedia('(prefers-color-scheme: dark)');

    const handleThemeChange = (e: MediaQueryListEvent) => {
      setTheme(e.matches);
      console.log('test color listener');
    };

    const editor = monaco.editor.create(container!, {
      value,
      language,
      lineNumbers: 'off',
      glyphMargin: false,
      folding: false,
      lineDecorationsWidth: 0,
      lineNumbersMinChars: 0,
      fontSize: 16,
      minimap: {
        enabled: false,
      },
    });
    editor.onDidChangeModelContent((event) => {
      value = editor.getValue();
    });

    // Initialize dark mode if set
    setTheme(mqList.matches);

    // Listen to system theme preference changes and set color theme
    mqList.addEventListener('change', handleThemeChange);

    return () => {
      mqList.removeEventListener('change', handleThemeChange);
      editor.dispose();
    };
  });
</script>

<div {id} bind:this={container} style={`height: ${height}`} />

<style>
  div {
    border-radius: 4px;
    overflow: hidden;
    padding-top: 6px;
    padding-left: 12px;
    box-shadow: var(--heavy-3d-box-shadow);
  }
</style>
