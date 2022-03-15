<script lang="ts">
  import { shortcuts } from '../../lib/actions/shortcut';

  export let filterString: string = '';
  export let clearEvents: () => void;
</script>

<div class="container">
  <input type="text" placeholder="Filter..." bind:value={filterString} />
  <button
    use:shortcuts={{
      chords: [
        {
          meta: true,
          code: 'KeyK',
        },
        {
          control: true,
          code: 'KeyL',
        },
      ],
    }}
    on:click={(e) => {
      clearEvents();
    }}
  >
    Clear Logs
  </button>
</div>

<style>
  .container {
    display: flex;
  }

  input {
    width: 100%;
    border: none;
    font-size: 16px;
    background: var(--primary-bg-color);
    padding: 8px 12px;
  }

  input,
  button {
    border: none;
    white-space: nowrap;
    border-top: 1px solid var(--layout-bg-color);
    outline: none;
  }

  button {
    font-size: 0.85em;
    padding: 8px 16px;
    border-left: 1px solid var(--layout-bg-color);
    background: var(--primary-bg-color);
  }

  button:hover {
    background: var(--grey-e-color);
  }

  button:active {
    background: var(--grey-c-color);
  }
</style>
