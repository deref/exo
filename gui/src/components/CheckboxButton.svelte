<script lang="ts">
  import CheckSVG from './mono/CheckSVG.svelte';

  export let active: boolean = false;
  export let tooltip: string = '';
</script>

<button
  on:click
  class={active ? 'active' : ''}
  title={tooltip !== '' ? tooltip : undefined}
>
  {#if active}
    <div>
      <CheckSVG />
    </div>
  {/if}
</button>

<style>
  button {
    font-family: inherit;
    font-size: 0;
    font-weight: 450;
    width: 22px;
    height: 22px;
    padding: 0;
    color: var(--checkbox-color);
    background-color: var(--checkbox-bg-color);
    border: 1px solid var(--checkbox-border-color);
    border-radius: 4px;
    outline: none;
    margin-left: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  button > div {
    width: 32px;
    height: 30px;
    padding-left: 4px;
    padding-bottom: 2px;
  }

  button > div > :global(svg) {
    width: 28px;
    height: 28px;
  }

  button.active {
    border-color: var(--checkbox-active-border-color);
  }

  button :global(svg),
  button :global(svg *) {
    fill: currentColor;
  }

  button.active :global(svg),
  button.active :global(svg *) {
    fill: var(--checkbox-active-color);
  }

  button:hover {
    background-color: var(--checkbox-hover-bg-color);
    color: var(--checkbox-hover-color);
  }

  button.active:hover {
    background-color: var(--checkbox-active-hover-bg-color);
  }

  button:focus {
    background-color: var(--checkbox-focus-bg-color);
  }

  button.active:focus {
    background-color: var(--checkbox-active-focus-bg-color);
  }
</style>
