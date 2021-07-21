<script lang="ts">
  import Textarea from './Textarea.svelte';

  export let name: string;
  export let id: string;
  export let value: string[] = [];
  
  let text = value.join('\n');;
  
  $: {
    value = text.split('\n')
      .map(line => line.trim())
      .filter(line => line);
  }
</script>

<Textarea
  bind:value={text}
  id={id}
  name={name}
/>
