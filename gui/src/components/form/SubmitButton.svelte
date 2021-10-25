<script lang="ts">
  import Button from '../Button.svelte';

  // This is a right-aligned submit button for forms.
  // Used for create/edit forms where the code below
  // would otherwise be highly repetitive.
</script>

<div><Button type="submit"><slot /></Button></div>

<style>
  div {
    display: flex;
    justify-content: flex-end;
    margin: 24px 0;
  }
</style>
