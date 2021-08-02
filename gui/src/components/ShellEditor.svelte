<script lang="ts">
  import TextEditor from './TextEditor.svelte';

  export let id: string | undefined;
  export let value: string = '';
</script>

<TextEditor {id} language="shell" bind:value />
