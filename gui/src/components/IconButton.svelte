<script lang="ts">
  export let active: boolean = false;
  export let tooltip: string = '';
</script>

<button
  on:click
  class={active ? 'active' : ''}
  title={tooltip !== '' ? tooltip : undefined}
>
  <slot />
</button>

<style>
  button {
    font-family: inherit;
    font-size: 0;
    font-weight: 450;
    padding: 6px;
    color: #777777;
    background-color: #77777700;
    border-radius: 4px;
    border: none;
    outline: none;
    margin-left: 4px;
  }

  button.active {
    background-color: #55ccff33;
  }

  button :global(svg),
  button :global(svg *) {
    fill: currentColor;
  }

  button:hover {
    background-color: #77777711;
    color: #444444;
  }

  button.active:hover {
    background-color: #55ccff44;
  }

  button:focus {
    background-color: #77777733;
  }

  button.active:focus {
    background-color: #55ccff66;
  }
</style>
