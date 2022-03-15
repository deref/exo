<script lang="ts">
  export let loading = false;
</script>

<div class="container">
  <div class="content" class:loading />
</div>

<style>
  .container {
    display: absolute;
    left: 0;
    top: 0;
    right: 0;
    height: 0px;
  }

  .content {
    width: 100%;
    height: 1px;
  }

  .loading {
    /* TODO: Use a distinct color variable. */
    background: var(--link-color);
    box-shadow: 0 0 8px var(--link-color);
    animation: pulse 3s infinite;
  }

  @keyframes pulse {
    0% {
      opacity: 0%;
    }
    50% {
      opacity: 100%;
    }
    100% {
      opacity: 0%;
    }
  }
</style>
