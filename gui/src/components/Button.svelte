<script lang="ts">

export let type: string | undefined = undefined;

</script>

<button on:click type={type}>
  <slot/>
</button>

<style>

button {
  background: linear-gradient(#fff, #f5f5f5);
  border: none;
  border-radius: 6px;
  padding: 12px 18px;
  box-shadow:
    0 .33px 0 1px hsla(0,0%,100%,.15),
    0 4px 8px -3px rgba(0,0,0,.15),
    0 .4px 0 .8px rgba(0,0,0,.25);
}

button:hover {
  background: linear-gradient(#fafafa, #e7e7e7);
  box-shadow:
    0 .33px 0 1px hsla(0,0%,100%,.15),
    0 6px 8px -4px rgba(0,0,0,.2),
    0 .4px 0 .8px rgba(0,0,0,.35);
}

button:active {
  background: linear-gradient(#f7f7f7, #e0e0e0);
  box-shadow:
    0 .33px 0 1px hsla(0,0%,100%,.15),
    0 4px 6px -3px rgba(0,0,0,.15),
    0 .4px 0 .8px rgba(0,0,0,.45);
}

</style>
