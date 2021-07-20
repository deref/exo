<script lang="ts">
  export let name: string;
  export let value: string[] = [];
  
  let text = value.join('\n');;
  
  $: {
    value = text.split('\n')
      .map(line => line.trim())
      .filter(line => line);
  }
</script>

<textarea
  bind:value={text}
  name={name}
/>
