<script lang="ts">
  interface EditorMode {
    name: string;
    id: string;
    onActivate?(): void;
  }

  export let mode: string;
  export let editorModes: EditorMode[];
</script>

<div>
  <span>Edit as:</span>
  {#each editorModes as editorMode}
    <button
      class:selected={mode === editorMode.id}
      on:click|preventDefault={() => {
        if (mode === editorMode.id) {
          return;
        }
        mode = editorMode.id;
        if (editorMode.onActivate) {
          editorMode.onActivate();
        }
      }}
    >
      {editorMode.name}
    </button>
  {/each}
</div>

<style>
  div {
    margin: 24px 0;
  }

  span {
    margin-right: 8px;
  }

  button {
    border: none;
    background: none;
    font-weight: 450;
    color: var(--grey-7-color);
  }

  button:hover {
    text-decoration: underline;
  }

  .selected {
    text-decoration: underline;
    font-weight: 450;
    color: var(--link-color);
  }
</style>
