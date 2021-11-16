<script lang="ts">
  export let id: string | undefined = undefined;
  export let name: string | undefined = undefined;

  export let value: string;

  export let options: string[];
</script>

<select {id} bind:value {name} on:blur on:focus>
  {#each options as option}
    <option value={option}>{option}</option>
  {/each}
</select>

<style>
  select {
    border: none;
    margin: 0;
    padding: 12px 18px;
    background: var(--primary-bg-color);
    color: var(--strong-color);
    box-shadow: var(--text-input-shadow);
    width: var(--input-width);
    height: 2.5rem;
    outline: none;
    width: 100%;
    border-radius: 4px;
  }

  select:focus,
  select:focus-within {
    box-shadow: var(--text-input-shadow-focus) !important;
  }
</style>
