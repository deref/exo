<script lang="ts">
  export let id: string | undefined = undefined;
  export let name: string | undefined = undefined;

  export let value: string;
</script>

<textarea {id} bind:value {name} on:blur on:focus on:input />

<style>
  textarea {
    border: none;
    border-radius: 4px;
    padding: 12px 18px;
    background: var(--primary-bg-color);
    color: var(--strong-color);
    box-shadow: var(--text-input-shadow);
    width: 100%;
    height: 120px;
    resize: vertical;
    outline: none;
  }

  textarea:focus,
  textarea:focus-within {
    box-shadow: var(--text-input-shadow-focus) !important;
  }
</style>
