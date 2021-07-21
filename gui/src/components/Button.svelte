<script type="ts">

export let type: string | undefined;

</script>

<button type={type}>
  <slot/>
</button>

<style>

button {
  background: #efefef;
  border: none;
  border-radius: 6px;
  padding: 12px 18px;
  box-shadow:
    0 .33px 0 1px hsla(0,0%,100%,.15),
    0 6px 9px -4px rgba(0,0,0,.3),
    0 .4px 0 .8px rgba(0,0,0,.2);
}

</style>
