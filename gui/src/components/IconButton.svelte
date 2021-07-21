<button on:click>
  <slot/>
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

button :global(svg), button :global(svg *) {
  fill: currentColor;
}

button:hover {
  background-color: #77777711;
  color: #444444;
}

button:focus {
  background-color: #77777733;
}

</style>
