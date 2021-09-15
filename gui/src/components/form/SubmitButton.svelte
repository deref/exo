<script lang="ts">
  import Button from '../Button.svelte';
</script>

<div class="buttons">
  <Button type="submit"><slot /></Button>
</div>

<style>
  .buttons {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin: 24px 0;
  }
</style>
