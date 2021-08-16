<script lang="ts">
  export let id: string | undefined = undefined;
  export let name: string | undefined = undefined;
  export let placeholder: string | undefined = undefined;

  export let value: string;
</script>

<input {id} bind:value {name} {placeholder} on:blur on:focus on:input />

<style>
  input {
    border: none;
    border-radius: 4px;
    padding: 12px 18px;
    box-shadow: var(--text-input-shadow);
    width: var(--input-width);
  }
</style>
