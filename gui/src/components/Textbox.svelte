<script lang="ts">
  export let id: string | undefined = undefined;
  export let name: string | undefined = undefined;
  export let placeholder: string | undefined = undefined;
  export let autofocus: boolean | undefined = undefined;

  export let value: string;
</script>

<input
  {id}
  bind:value
  {name}
  {placeholder}
  on:blur
  on:focus
  on:input
  {autofocus}
/>

<style>
  input {
    border: none;
    border-radius: 4px;
    padding: 12px 18px;
    background: var(--primary-bg-color);
    color: var(--strong-color);
    box-shadow: var(--text-input-shadow);
    width: var(--input-width);
    outline: none;
  }

  input:focus,
  input:focus-within {
    box-shadow: var(--text-input-shadow-focus);
  }
</style>
