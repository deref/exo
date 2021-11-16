<script lang="ts">
  export let id: string | undefined = undefined;
  export let name: string | undefined = undefined;
  export let placeholder: string | undefined = undefined;

  export let value: string;

  const matchGroup = (index: 1 | 2) => {
    if (value) {
      const match = value.match(/([0-9]+)(\w+)/);

      if (match) {
        return match[index];
      }
    }

    return '';
  };

  let num = matchGroup(1);
  let unit = matchGroup(2);

  export let unitOptions: string[] | undefined = undefined;

  $: {
    value = num + unit;
  }
</script>

<input bind:value style="display:none" />

<div>
  <input
    type="number"
    {id}
    bind:value={num}
    {name}
    {placeholder}
    on:blur
    on:focus
    on:input
  />

  {#if unitOptions}
    <select bind:value={unit}>
      {#each unitOptions as unitOption}
        <option value={unitOption}>{unitOption}</option>
      {/each}
    </select>
  {/if}
</div>

<style>
  input,
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
  }

  div {
    display: grid;
    grid-template-columns: auto max-content;
    gap: 1px;
  }

  div > *:first-child {
    border-top-left-radius: 4px;
    border-bottom-left-radius: 4px;
  }

  div > *:last-child {
    border-top-right-radius: 4px;
    border-bottom-right-radius: 4px;
  }

  input:focus,
  input:focus-within,
  select:focus,
  select:focus-within {
    box-shadow: var(--text-input-shadow-focus) !important;
  }
</style>
