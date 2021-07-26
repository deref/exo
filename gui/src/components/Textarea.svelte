<script lang="ts">
  export let id: string | undefined = undefined;
  export let name: string | undefined = undefined;

  export let value: string;
</script>

<textarea {id} bind:value {name} on:blur on:focus on:input />

<style>
  textarea {
    border: none;
    border-radius: 6px;
    padding: 12px 18px;
    box-shadow: 0 0.33px 0 1px hsla(0, 0%, 100%, 0.15),
      0 6px 9px -4px rgba(0, 0, 0, 0.1) inset,
      0 0.4px 0 0.8px rgba(0, 0, 0, 0.1) inset;
    width: 100%;
    height: 120px;
    resize: vertical;
  }
</style>
